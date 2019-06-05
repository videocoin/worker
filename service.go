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

	profiles_v1 "github.com/VideoCoin/cloud-api/profiles/v1"
	transcoder_v1 "github.com/VideoCoin/cloud-api/transcoder/v1"
	workorder_v1 "github.com/VideoCoin/cloud-api/workorder/v1"
	bc "github.com/VideoCoin/cloud-pkg/bcops"
	"github.com/VideoCoin/cloud-pkg/stream"
	"github.com/VideoCoin/cloud-pkg/streamManager"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/denisbrodbeck/machineid"
	"github.com/nats-io/go-nats"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
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

	key, err := bc.LoadBcPrivKeys(cfg.KeyFile, cfg.Password, bc.FromMemory)
	if err != nil {
		log.Fatalf("failed to load private keys: %s", err.Error())
	}

	bcAuth, err := bc.GetBCAuth(client, key)
	if err != nil {
		log.Fatalf("failed to get blockchain auth: %s", err.Error())
	}

	nc, err := nats.Connect(cfg.NatsURL)
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

	err = s.registerTranscoder(&transcoder_v1.Transcoder{
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

func (s *Service) handleTranscode(workOrder *workorder_v1.WorkOrder, profile *profiles_v1.Profile, uid string) error {
	s.log.Infof("transcoding: %d using input: %s with stream_id: %d", workOrder.Id, workOrder.TranscodeInputUrl, workOrder.StreamId)

	dir := path.Join(s.cfg.OutputDir, workOrder.StreamHash)
	m3u8 := path.Join(dir, "index.m3u8")

	cmd := buildCmd(workOrder.TranscodeInputUrl, dir, profile)
	var stopChan = make(chan struct{})

	go s.handleStop(uid, stopChan, cmd)
	s.listenForStop(uid, stopChan)

	fullDir := fmt.Sprintf("%s/%d", dir, profile.Bitrate)
	err := prepareDir(fullDir)

	if err != nil {
		return err
	}

	go s.syncDir(stopChan, cmd, workOrder, fullDir, profile.Bitrate)

	if err := s.generatePlaylist(workOrder.StreamHash, m3u8, profile.Bitrate); err != nil {
		return err
	}

	go s.transcode(cmd,
		stopChan,
		workOrder.TranscodeInputUrl,
		uid,
	)

	return nil
}

func (s *Service) transcode(
	cmd *exec.Cmd,
	stop chan struct{},
	streamurl string,
	uid string,
) {
	s.waitForStreamReady(streamurl)

	_, err := cmd.CombinedOutput()
	if err != nil {
		s.log.Errorf("failed to transcode: err : %s cmd args [ %v ]", err.Error(), cmd.Args)
	}

	stop <- struct{}{}

	s.log.Info("transcode complete")
}

func (s *Service) handleStop(uid string, stopChan chan struct{}, cmd *exec.Cmd) {
	<-stopChan
	cmd.Process.Signal(os.Interrupt)
	if err := s.updateTranscoderStatus(uid, transcoder_v1.TranscoderStatusAvailable); err != nil {
		s.log.Warnf("failed to update transcode status: %s", err.Error())
	}
	close(stopChan)
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
	args := fmt.Sprintf("-live_start_index 0 -b:v %d -vf scale=%d:-2 -strict -2 -c:v libx264 -c:a aac -r %f -bsf:v h264_mp4toannexb -map 0 -f segment -segment_time 2 -segment_format mpegts -segment_list %s/%d/index.m3u8 -segment_list_type m3u8 %s/%d/%%d.ts", profile.Bitrate, profile.Width, profile.Fps, dir, profile.Bitrate, dir, profile.Bitrate)
	process = append(process, strings.Split(args, " ")...)

	return exec.Command("ffmpeg", process...)
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
