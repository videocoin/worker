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

	"github.com/denisbrodbeck/machineid"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.WithField("service", "transcode")
	cfg := LoadConfig()

	level, _ := logrus.ParseLevel(cfg.LogLevel)
	logrus.SetLevel(level)
	// Generate unique connection name
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}

	managerConn, err := grpc.Dial(cfg.ManagerRPCADDR, opts...)
	if err != nil {
		log.Fatalf("failed to dial manager: %s", err.Error())
	}

	verifierConn, err := grpc.Dial(cfg.VerifierRPCADDR, opts...)
	if err != nil {
		log.Fatalf("failed to dial verifier: %s", err.Error())
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
		log.Fatalf("failed to dial blockchain: %s", err.Error())
	}

	managerAddress := common.HexToAddress(cfg.SMCA)

	sm, err := streamManager.NewManager(managerAddress, client)
	if err != nil {
		log.Fatalf("failed to make new stream manager: %s", err.Error())
	}

	key, err := bc.LoadBcPrivKeys(cfg.KeyFile, cfg.Password)
	if err != nil {
		log.Fatalf("failed to load private keys: %s", err.Error())
	}

	bcAuth, err := bc.GetBCAuth(client, key)
	if err != nil {
		log.Fatalf("failed to get blockchain auth: %s", err.Error())
	}

	return &Service{
		streamManager: sm,
		bcAuth:        bcAuth,
		bcClient:      client,
		cfg:           cfg,
		manager:       manager,
		verifier:      v,
		ctx:           ctx,
		log:           log,
	}, nil

}

// Start creates new service and blocks until stop signal
func Start() error {
	s, err := New()
	if err != nil {
		s.log.Errorf("failed to create service: %s", err.Error())
		return err
	}

	workOrder, err := s.manager.GetJob(s.ctx, &pb.GetJobRequest{})
	if err != nil {
		s.log.Debugf("failed to get job: %s", err.Error())
		return err
	}

	profile, err := s.manager.GetProfile(s.ctx, &pb.GetProfileRequest{ProfileId: workOrder.Profile})
	if err != nil {
		s.log.Debugf("failed to get profile: %s", err.Error())
		return err
	}

	streamInstance, err := stream.NewStream(common.HexToAddress(workOrder.ContractAddress), s.bcClient)
	if err != nil {
		s.log.Errorf("failed to create new stream: %s", err.Error())
		return err
	}

	s.streamInstance = streamInstance

	if err = s.handleTranscodeTask(workOrder, profile); err != nil {
		s.log.Errorf("failed to handle transcode task: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) register() {
	info, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	machineID, _ := machineid.ProtectedID(s.cfg.HashKey)

	s.manager.RegisterTranscoder(context.Background(), &pb.Transcoder{
		Id:          machineID,
		CpuCores:    info[0].Cores,
		CpuMhz:      info[0].Mhz,
		TotalMemory: memInfo.Total,
	})
}

func (s *Service) handleTranscodeTask(workOrder *pb.WorkOrder, profile *pb.Profile) error {

	s.log.Infof("starting transcode task: %d using input: %s with stream_id: %d", workOrder.Id, workOrder.InputUrl, workOrder.StreamId)

	dir := path.Join(s.cfg.OutputDir, fmt.Sprintf("%d", workOrder.StreamId))
	m3u8 := path.Join(dir, "index.m3u8")

	cmd := buildCmd(workOrder.InputUrl, dir, profile)
	var stopChan = make(chan struct{})

	for _, b := range bitrates {

		fullDir := fmt.Sprintf("%s/%d", dir, b)
		err := prepareDir(fullDir)

		if err != nil {
			return err
		}
		go s.SyncDir(stopChan, cmd, workOrder, fullDir, b)

	}

	if err := s.GeneratePlaylist(workOrder.StreamId, m3u8); err != nil {
		return err
	}

	s.transcode(cmd, stopChan, workOrder.InputUrl, workOrder.StreamId)

	return nil
}

func (s *Service) transcode(cmd *exec.Cmd, stop chan struct{}, streamurl string, streamID int64) {
	s.waitForStreamReady(streamurl)
	s.log.Info("starting transcode")
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.log.Errorf("failed to transcode: err : %s output: %s", err.Error(), string(out))
	}

	stop <- struct{}{}
	s.log.Info("calling refund")
	if err := s.refund(streamID); err != nil {
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
		time.Sleep(10 * time.Second)
	}
}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func buildCmd(inputURL string, dir string, profile *pb.Profile) *exec.Cmd {
	process := []string{"-re", "-i", inputURL}

	for _, b := range bitrates {
		args := fmt.Sprintf("-live_start_index 0 -b:v %d -vf scale=%d:-2 -strict -2 -c:v libx264 -c:a aac -r %f -bsf:v h264_mp4toannexb -map 0 -f segment -segment_time 10 -segment_format mpegts -segment_list %s/%d/index.m3u8 -segment_list_type m3u8 %s/%d/%%d.ts", profile.Bitrate, profile.Width, profile.Fps, dir, b, dir, b)
		process = append(process, strings.Split(args, " ")...)

	}

	cmd := exec.Command("ffmpeg", process...)

	return cmd
}

func (s *Service) refund(streamID int64) error {
	_, err := s.streamInstance.Refund(s.bcAuth)
	if err != nil {
		return err
	}

	_, err = s.manager.UpdateStreamStatus(s.ctx, &pb.UpdateStreamStatusRequest{StreamId: streamID, Status: pb.WorkOrderStatusCompleted.String(), Refunded: true})
	if err != nil {
		s.log.Warnf("failed to update stream status: %s", err.Error())
	}

	return nil
}
