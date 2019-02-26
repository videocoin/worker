package transcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	bc "github.com/VideoCoin/common/bcops"
	pb "github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

// New initialize and return a new Service object
func New() (*Service, error) {
	cfg := LoadConfig()

	// Generate unique connection name
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}

	managerConn, err := grpc.Dial(cfg.ManagerRPCADDR, opts...)
	if err != nil {
		panic(err)
	}

	verifierConn, err := grpc.Dial(cfg.VerifierRPCADDR, opts...)
	if err != nil {
		panic(err)
	}

	manager := pb.NewManagerServiceClient(managerConn)
	status, err := manager.Health(context.Background(), &empty.Empty{})
	if status.GetStatus() != "healthy" || err != nil {
		return nil, fmt.Errorf("failed to get healthy status from manager: %v", err)
	}

	v := pb.NewVerifierServiceClient(verifierConn)
	status, err = v.Health(context.Background(), &empty.Empty{})
	if status.GetStatus() != "healthy" || err != nil {
		return nil, fmt.Errorf("failed to get healthy status from verifier: %v", err)
	}

	ctx := context.Background()

	client, err := ethclient.Dial(cfg.BlockchainURL)
	if err != nil {
		panic(err)
	}

	managerAddress := common.HexToAddress(cfg.SMCA)

	sm, err := streamManager.NewManager(managerAddress, client)
	if err != nil {
		panic(err)
	}

	key, err := bc.LoadBcPrivKeys(cfg.KeyFile, cfg.Password)
	if err != nil {
		panic(err)
	}

	bcAuth, err := bc.GetBCAuth(client, key)
	if err != nil {
		panic(err)
	}

	return &Service{
		streamManager: sm,
		bcAuth:        bcAuth,
		bcClient:      client,
		cfg:           cfg,
		manager:       manager,
		verifier:      v,
		ctx:           ctx,
		log:           logrus.WithField("name", "xcode"),
	}, nil

}

// Start creates new service and blocks until stop signal
func Start() error {
	s, err := New()
	if err != nil {
		panic(err)
	}

	workOrder, err := s.manager.GetJob(s.ctx, &pb.GetJobRequest{})
	if err != nil {
		return err
	}

	profile, err := s.manager.GetProfile(s.ctx, &pb.GetProfileRequest{ProfileId: workOrder.Profile})
	if err != nil {
		return err
	}

	streamInstance, err := stream.NewStream(common.HexToAddress(workOrder.WalletAddress), s.bcClient)
	if err != nil {
		return err
	}

	s.streamInstance = streamInstance

	if err = s.handleTranscodeTask(workOrder, profile); err != nil {
		s.log.Errorf("failed to handle transcode task: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) handleTranscodeTask(workOrder *pb.WorkOrder, profile *pb.Profile) error {

	s.log.Infof("starting transcode task: %d using input: %s with stream_id: %d", workOrder.Id, workOrder.InputUrl, workOrder.StreamId)

	dir := path.Join(s.cfg.OutputDir, fmt.Sprintf("%d", workOrder.StreamId))
	m3u8 := path.Join(dir, "index.m3u8")

	var stopChan = make(chan bool)
	for _, b := range bitrates {

		fullDir := fmt.Sprintf("%s/%d", dir, b)
		err := prepareDir(fullDir)

		if err != nil {
			return err
		}
		go s.SyncDir(stopChan, workOrder, fullDir, b)

	}

	if err := s.GeneratePlaylist(workOrder.StreamId, m3u8); err != nil {
		return err
	}

	cmd := buildCmd(workOrder.InputUrl, dir, profile)

	go s.monitorBalance(cmd, stopChan, workOrder.ContractAddress)

	s.transcode(cmd, stopChan, workOrder.InputUrl)

	return nil
}

func (s *Service) monitorBalance(cmd *exec.Cmd, stop chan bool, addr string) {
	for {
		time.Sleep(10 * time.Second)
		balance, err := s.manager.CheckBalance(context.Background(), &pb.CheckBalanceRequest{ContractAddress: addr})
		if err != nil {
			s.log.Warnf("failed to check balance, allowing work")
		}

		s.log.Infof("current balance at address %s is %d", addr, balance.Balance)

		if balance.Balance <= 0 {
			cmd.Process.Kill()
			stop <- true
		}
	}
}

func (s *Service) transcode(cmd *exec.Cmd, stop chan bool, streamurl string) {
	s.waitForStreamReady(streamurl)
	s.log.Info("starting transcode")
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.log.Errorf("failed to transcode: err : %s output: %s", err.Error(), string(out))
	}

	stop <- true
	if err := s.refund(); err != nil {
		s.log.Errorf("failed to refund:%s", err.Error())
	}

	s.log.Info("transcode complete")
}

func (s *Service) waitForStreamReady(streamurl string) {
	maxretry := 10
	for i := 0; i < maxretry; i++ {
		resp, _ := http.Head(streamurl)
		if resp.StatusCode == 200 {
			return
		}
		s.log.Infof("waiting for stream %s to become ready...", streamurl)
		time.Sleep(30 * time.Second)
	}
}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func buildCmd(inputURL string, dir string, profile *pb.Profile) *exec.Cmd {
	process := []string{"-re", "-i", inputURL}

	for _, b := range bitrates {
		args := fmt.Sprintf("-live_start_index 0 -b:v %d -vf scale=%d:-2 -strict -2 -r %f -codec copy -bsf:v h264_mp4toannexb -map 0 -f segment -segment_time 10 -segment_format mpegts -segment_list %s/%d/index.m3u8 -segment_list_type m3u8 %s/%d/%%d.ts", profile.Bitrate, profile.Width, profile.Fps, dir, b, dir, b)
		process = append(process, strings.Split(args, " ")...)

	}

	cmd := exec.Command("ffmpeg", process...)

	return cmd
}

func (s *Service) refund() error {
	_, err := s.streamInstance.Refund(s.bcAuth)
	if err != nil {
		return err
	}
	return nil
}
