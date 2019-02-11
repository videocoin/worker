package transcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"log"

	bc "github.com/VideoCoin/common/bcops"
	pb "github.com/VideoCoin/common/proto"
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

	return s.handleTranscodeTask(workOrder)

}

func (s *Service) handleTranscodeTask(workOrder *pb.WorkOrder) error {

	s.log.Infof("starting transcode task: %d using input: %s with stream_id: %d", workOrder.Id, workOrder.InputUrl, workOrder.StreamId)

	dir := path.Join(s.cfg.OutputDir, fmt.Sprintf("%d", workOrder.StreamId))
	m3u8 := path.Join(dir, "index.m3u8")

	for _, b := range bitrates {

		fullDir := fmt.Sprintf("%s/%d", dir, b)
		err := prepareDir(fullDir)
		if err != nil {
			s.log.Errorf("failed to prepare directory [ %s ]: %s", fullDir, err.Error())
		}

		s.log.Infof("monitoring chunks in %s", fullDir)

		go s.monitorChunks(fullDir, workOrder)
		go s.SyncDir(workOrder, fullDir, b)

	}

	if err := s.GeneratePlaylist(workOrder.StreamId, m3u8); err != nil {
		s.log.Fatalf("failed to generate playlist: %s", err.Error())
	}

	cmd := buildCmd(workOrder.InputUrl, dir)

	go s.monitorBalance(cmd, workOrder.ContractAddress)

	s.transcode(cmd, workOrder.InputUrl)

	return nil
}

func (s *Service) monitorBalance(cmd *exec.Cmd, addr string) {
	for {
		time.Sleep(10 * time.Second)
		balance, err := s.manager.CheckBalance(context.Background(), &pb.CheckBalanceRequest{ContractAddress: addr})
		if err != nil {
			s.log.Warnf("failed to check balance, allowing work")
		}

		if balance.Balance <= 0 {
			cmd.Process.Kill()
		}
	}
}

func (s *Service) monitorChunks(dir string, task *pb.WorkOrder) {
	for {
		time.Sleep(2 * time.Second)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Printf("failed to read dir: %s", err.Error())
		}

		if len(files) < 2 {
			continue
		}

		break
	}

	task.Status = pb.WorkOrderStatusReady.String()
	update := &pb.UpdateStreamStatusRequest{
		StreamId: task.StreamId,
		Status:   task.Status,
	}
	_, err := s.manager.UpdateStreamStatus(s.ctx, update)

	if err != nil {
		s.log.Errorf("failed to update stream status: %s", err.Error())
	}
	fmt.Println(task.Status, task.StreamId)

}

func (s *Service) transcode(cmd *exec.Cmd, streamurl string) {
	waitForStreamReady(streamurl)
	log.Println("starting transcode")
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.log.Fatalf("failed to transcode: err : %s output: %s", err.Error(), string(out))
	}

	s.log.Infof("transcode complete")
}

func waitForStreamReady(streamurl string) {
	maxretry := 10
	for i := 0; i < maxretry; i++ {
		resp, _ := http.Head(streamurl)
		if resp.StatusCode == 200 {
			return
		}
		log.Printf("waiting for stream %s to become ready...", streamurl)
		time.Sleep(30 * time.Second)
	}
}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func buildCmd(inputURL string, dir string) *exec.Cmd {
	process := []string{"-re", "-i", inputURL}

	for _, b := range bitrates {
		args := fmt.Sprintf("-hls_allow_cache 0 -hls_flags append_list -f ssegment -b:v %d -strict -2 -c:v h264 -profile:v main -segment_list_flags live -segment_time 10 -segment_format mpegts -an -segment_list %s/%d/index.m3u8 %s/%d/%%d.ts", b, dir, b, dir, b)
		process = append(process, strings.Split(args, " ")...)

	}

	cmd := exec.Command("ffmpeg", process...)

	return cmd

}
