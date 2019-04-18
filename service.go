package transcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	manager_v1 "github.com/VideoCoin/cloud-api/manager/v1"
	profiles_v1 "github.com/VideoCoin/cloud-api/profiles/v1"
	transcoder_v1 "github.com/VideoCoin/cloud-api/transcoder/v1"
	verifier_v1 "github.com/VideoCoin/cloud-api/verifier/v1"
	workorder_v1 "github.com/VideoCoin/cloud-api/workorder/v1"
	bc "github.com/VideoCoin/common/bcops"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/denisbrodbeck/machineid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nats-io/go-nats"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// New initialize and return a new Service object
func newService() (*Service, error) {
	log := logrus.WithField("service", "transcode")
	cfg := LoadConfig()

	level, _ := logrus.ParseLevel(cfg.LogLevel)
	logrus.SetLevel(level)
	// Generate unique connection name
	b := make([]byte, 16)
	rand.Read(b)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	managerConn, err := grpc.Dial(cfg.ManagerRPCADDR, opts...)
	if err != nil {
		log.Fatalf("failed to dial manager: %s", err.Error())
	}

	verifierConn, err := grpc.Dial(cfg.VerifierRPCADDR, opts...)
	if err != nil {
		log.Fatalf("failed to dial verifier: %s", err.Error())
	}

	manager := manager_v1.NewManagerServiceClient(managerConn)
	status, err := manager.Health(context.Background(), &empty.Empty{})
	if status.GetStatus() != "healthy" || err != nil {
		return nil, fmt.Errorf("failed to get healthy status from manager")
	}

	v := verifier_v1.NewVerifierServiceClient(verifierConn)
	status, err = v.Health(context.Background(), &empty.Empty{})
	if status.GetStatus() != "healthy" || err != nil {
		return nil, fmt.Errorf("failed to get healthy status from verifier")
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

	nc, err := nats.Connect(cfg.NatsURL, nats.Token(cfg.NatsToken))
	if err != nil {
		log.Fatalf("failed to connect to nats: %s", err.Error())
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("failed to creat encoded connection: %s", err.Error())
	}

	return &Service{
		streamManager: sm,
		bcAuth:        bcAuth,
		bcClient:      client,
		cfg:           cfg,
		manager:       manager,
		verifier:      v,
		ec:            ec,
		nc:            nc,
		ctx:           ctx,
		log:           log,
	}, nil

}

// Start creates new service and blocks until stop signal
func Start() error {
	s, err := newService()
	if err != nil {
		return err
	}

	uid, err := machineid.ProtectedID(cfg.HashKey)
	if err != nil {
		s.log.Warnf("failed to calculate machine id: %s", err.Error())
	}

	{

		go s.subscribe(uid)
		s.register(uid)
	}

	s.wait()

	return nil
}

func (s *Service) register(uid string) {
	info, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()

	s.manager.RegisterTranscoder(context.Background(), &transcoder_v1.Transcoder{
		Id:          uid,
		CpuCores:    info[0].Cores,
		CpuMhz:      info[0].Mhz,
		TotalMemory: memInfo.Total,
		Status:      transcoder_v1.TranscoderStatusAvailable.String(),
	})
}

func (s *Service) handleTranscode(workOrder *workorder_v1.WorkOrder, profile *profiles_v1.Profile, uid string) error {
	s.log.Infof("starting transcode: %d using input: %s with stream_id: %d",
		workOrder.Id, workOrder.InputUrl, workOrder.StreamId,
	)

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

		go s.syncDir(stopChan, cmd, workOrder, fullDir, b)

	}

	if err := s.generatePlaylist(workOrder.StreamId, m3u8); err != nil {
		return err
	}

	go s.transcode(cmd,
		stopChan,
		workOrder.InputUrl,
		workOrder.StreamId,
		workOrder.ContractAddress,
		uid,
	)

	return nil

}

func (s *Service) transcode(
	cmd *exec.Cmd,
	stop chan struct{},
	streamurl string,
	streamID int64,
	contractAddr string,
	uid string,
) {
	s.waitForStreamReady(streamurl)
	s.log.Info("starting transcode")
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.log.Errorf("failed to transcode: err : %s output: %s",
			err.Error(), string(out),
		)
	}

	stop <- struct{}{}
	s.log.Info("calling refund")
	if err := s.refund(streamID, contractAddr); err != nil {
		s.log.Errorf("failed to refund:%s", err.Error())
	}

	s.manager.UpdateTranscoderStatus(s.ctx, &manager_v1.UpdateTranscoderStatusRequest{TranscoderId: uid, Status: transcoder_v1.TranscoderStatusAvailable.String()})

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

func buildCmd(inputURL string, dir string, profile *profiles_v1.Profile) *exec.Cmd {
	process := []string{"-re", "-i", inputURL}

	for _, b := range bitrates {
		args := fmt.Sprintf("-live_start_index 0 -b:v %d -vf scale=%d:-2 -strict -2 -c:v libx264 -c:a aac -r %f -bsf:v h264_mp4toannexb -map 0 -f segment -segment_time 10 -segment_format mpegts -segment_list %s/%d/index.m3u8 -segment_list_type m3u8 %s/%d/%%d.ts", profile.Bitrate, profile.Width, profile.Fps, dir, b, dir, b)
		process = append(process, strings.Split(args, " ")...)
	}

	cmd := exec.Command("ffmpeg", process...)

	return cmd

}

func (s *Service) refund(streamID int64, addr string) error {
	streamInstance, err := s.createStreamInstance(addr)
	if err != nil {
		return err
	}

	_, err = streamInstance.Refund(s.bcAuth)
	if err != nil {
		return err
	}

	_, err = s.manager.UpdateStreamStatus(s.ctx, &manager_v1.UpdateStreamStatusRequest{
		StreamId: streamID,
		Status:   workorder_v1.WorkOrderStatusCompleted.String(),
		Refunded: true,
	})

	if err != nil {
		s.log.Warnf("failed to update stream status: %s", err.Error())
	}

	return nil

}

func (s *Service) createStreamInstance(addr string) (*stream.Stream, error) {
	streamInstance, err := stream.NewStream(common.HexToAddress(addr), s.bcClient)
	if err != nil {
		s.log.Errorf("failed to create new stream instance: %s", err.Error())
		return nil, err
	}

	return streamInstance, nil
}

func (s *Service) wait() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	s.log.Info("shutting down")
}
