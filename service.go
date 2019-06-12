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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogo/protobuf/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	manager_v1 "github.com/videocoin/cloud-api/manager/v1"
	profiles_v1 "github.com/videocoin/cloud-api/profiles/v1"
	transcoder_v1 "github.com/videocoin/cloud-api/transcoder/v1"
	verifier_v1 "github.com/videocoin/cloud-api/verifier/v1"
	bc "github.com/videocoin/cloud-pkg/bcops"
	"github.com/videocoin/cloud-pkg/stream"
	"github.com/videocoin/cloud-pkg/streamManager"
	"github.com/videocoin/cloud-pkg/uuid4"
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
	b := make([]byte, 8)
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
	status, err := manager.Health(context.Background(), &types.Empty{})
	if status.GetStatus() != "healthy" || err != nil {
		return nil, fmt.Errorf("failed to get healthy status from manager")
	}

	v := verifier_v1.NewVerifierServiceClient(verifierConn)
	status, err = v.Health(context.Background(), &types.Empty{})
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

	key, err := bc.LoadBcPrivKeys(cfg.Key, cfg.Secret, bc.FromMemory)
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
	s, err := newService()
	if err != nil {
		return err
	}

	uid, err := uuid4.New()
	if err != nil {
		return err
	}

	go s.wait()

	s.register(uid)

	s.pollForWork(uid)

	return nil
}

func (s *Service) register(uid string) {

	var (
		cores    int32
		mhz      float64
		memtotal uint64
	)

	info, err := cpu.Info()
	if err != nil {
		s.log.Errorf("failed to get cpu info: %s", err.Error())
	} else {
		cores = info[0].Cores
		mhz = info[0].Mhz
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		s.log.Errorf("failed to get memory info: %s", err.Error())
	} else {
		memtotal = memInfo.Total
	}

	_, err = s.manager.RegisterTranscoder(context.Background(), &transcoder_v1.Transcoder{
		Id:          uid,
		CpuCores:    cores,
		CpuMhz:      mhz,
		TotalMemory: memtotal,
		Status:      transcoder_v1.TranscoderStatusAvailable,
	})

	if err != nil {
		s.log.Errorf("failed to register transcoder: %s", err.Error())
	}
}

func (s *Service) pollForWork(uid string) {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		assignment, err := s.manager.GetWork(context.Background(), &types.Empty{})
		if assignment == nil {
			s.log.Debug("no work available")
			continue
		} else if err != nil {
			s.log.Errorf("failed to get work: %s", err.Error())
			continue
		}

		s.log.Info("work found")

		s.handleTranscode(assignment, uid)
	}
}

func (s *Service) handleTranscode(a *transcoder_v1.Assignment, uid string) error {
	s.log.Infof("transcoding: %d\nusing input: %s\nwith stream_id: %d", a.Workorder.Id, a.Workorder.TranscodeInputUrl, a.Workorder.StreamId)

	dir := path.Join(s.cfg.OutputDir, a.Workorder.StreamHash)
	m3u8 := path.Join(dir, "index.m3u8")

	cmd := buildCmd(a.Workorder.TranscodeInputUrl, dir, a.Profile)
	var stopChan = make(chan struct{})

	fullDir := fmt.Sprintf("%s/%d", dir, a.Profile.Bitrate)
	err := prepareDir(fullDir)

	if err != nil {
		return err
	}

	go s.syncDir(stopChan, cmd, a.Workorder, fullDir, a.Profile.Bitrate)

	if err := s.generatePlaylist(a.Workorder.StreamHash, m3u8, a.Profile.Bitrate); err != nil {
		return err
	}

	go s.transcode(cmd,
		stopChan,
		a.Workorder.TranscodeInputUrl,
		a.Workorder.StreamAddress,
		a.Workorder.StreamHash,
		uid,
	)

	return nil
}

func (s *Service) transcode(
	cmd *exec.Cmd,
	stop chan struct{},
	streamurl string,
	contractAddr string,
	streamHash string,
	uid string,
) {
	s.waitForStreamReady(streamurl)

	out, err := cmd.CombinedOutput()
	if err != nil {
		s.log.Errorf("failed to transcode: err : %s output: %s",
			err.Error(), string(out),
		)
	}

	stop <- struct{}{}

	s.manager.UpdateTranscoderStatus(s.ctx, &manager_v1.TranscoderStatusRequest{TranscoderId: uid, Status: transcoder_v1.TranscoderStatusAvailable})

	s.log.Info("transcode complete")

}

func (s *Service) waitForStreamReady(streamurl string) {
	maxretry := 15
	for i := 0; i < maxretry; i++ {
		resp, err := http.Head(streamurl)
		if err != nil {
			s.log.Errorf("failed to request stream: %s", err.Error())
		} else if resp.StatusCode == 200 {
			return
		}
		s.log.Infof("waiting for stream %s to become ready...", streamurl)
		time.Sleep(5 * time.Second)
	}

}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func buildCmd(inputURL string, dir string, profile *profiles_v1.Profile) *exec.Cmd {
	process := []string{"-re", "-i", inputURL}

	args := fmt.Sprintf("-live_start_index 0 -b:v %d -vf scale=%d:-2 -strict -2 -c:v libx264 -c:a aac -bsf:v h264_mp4toannexb -map 0 -f segment -segment_time 2 -segment_format mpegts -segment_list %s/%d/index.m3u8 -segment_list_type m3u8 %s/%d/%%d.ts", profile.Bitrate, profile.Width, dir, profile.Bitrate, dir, profile.Bitrate)
	process = append(process, strings.Split(args, " ")...)

	cmd := exec.Command("ffmpeg", process...)

	return cmd

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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	s.log.Info("shutting down")
}
