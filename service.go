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

	pb "github.com/VideoCoin/common/proto"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Byte constants
const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

var (
	bitrates = []uint32{
		2 * MB,
		4 * MB,
		8 * MB,
	}
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
		return nil, err
	}

	manager := pb.NewManagerServiceClient(managerConn)

	status, err := manager.Health(context.Background(), &empty.Empty{})
	if status.GetStatus() != "healthy" || err != nil {
		return nil, fmt.Errorf("failed to get healthy status: %v", err)
	}

	ctx := context.Background()

	return &Service{
		cfg:     cfg,
		manager: manager,
		ctx:     ctx,
		csyc:    CSyncInit(cfg, manager),
	}, nil

}

func (s *Service) reportStatus(userID string, streamID string, status string) error {
	ctx := context.Background()
	_, err := s.manager.UpdateStreamStatus(ctx, &pb.UpdateStreamStatusRequest{
		StreamId: streamID,
		Status:   status,
	})

	if err != nil {
		return err
	}

	return nil
}

// Start creates new service and blocks until stop signal
func Start() error {
	s, err := New()
	if err != nil {
		return err
	}

	task, err := s.manager.GetJob(s.ctx, &pb.GetJobRequest{})
	if err != nil {
		return err
	}

	task.Status = pb.WorkOrderStatusTranscoding.String()

	if _, err := s.manager.UpdateStreamStatus(s.ctx, &pb.UpdateStreamStatusRequest{
		StreamId: task.StreamId,
		Status:   task.Status,
	}); err != nil {
		log.Errorf("failed to report status")
	}

	return s.handleTranscodeTask(task)

}

func (s *Service) handleTranscodeTask(workOrder *pb.WorkOrder) error {

	log.Infof("starting transcode task: %d using input: %s", workOrder.Id, workOrder.InputUrl)

	dir := path.Join(s.cfg.OutputDir, workOrder.StreamId)
	m3u8 := path.Join(dir, "index.m3u8")

	for _, b := range bitrates {
		fullDir := fmt.Sprintf("%s/%d", dir, b)
		if err := prepareDir(fullDir); err != nil {
			log.Error(err.Error())
		}
		log.Infof("monitoring chunks in %s", fullDir)
		go s.monitorChunks(fullDir, workOrder)
		go s.csyc.SyncDir(workOrder, fullDir, b)
	}

	if err := generatePlaylist(m3u8); err != nil {
		panic(err)
	}

	args := buildCmd(workOrder.InputUrl, dir)

	transcode(args, workOrder.InputUrl)

	return nil
}

func (s *Service) monitorChunks(dir string, task *pb.WorkOrder) {
	for {
		time.Sleep(2 * time.Second)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Warnf("failed to read dir: %s", err.Error())
		}

		if len(files) < 2 {
			continue
		}

		break
	}

	task.Status = pb.WorkOrderStatusReady.String()

	if _, err := s.manager.UpdateStreamStatus(s.ctx, &pb.UpdateStreamStatusRequest{
		StreamId: task.StreamId,
		Status:   task.Status,
	}); err != nil {
		log.Errorf("failed to report status")
	}
}

func transcode(args []string, streamurl string) {
	waitForStreamReady(streamurl)
	log.Info("starting transcode")
	out, err := exec.Command("ffmpeg", args...).CombinedOutput()
	if err != nil {
		log.Errorf("failed to exec - output: %s", string(out))
		panic(err)
	}
	log.Info("transcode complete")
}

func generatePlaylist(filename string) error {
	m3u8 := []byte(`#EXTM3U
#EXT-X-VERSION:6
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1048576,RESOLUTION=640x360,CODECS="avc1.42e00a,mp4a.40.2"
360p/index.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=3145728,RESOLUTION=842x480,CODECS="avc1.42e00a,mp4a.40.2"
480p/index.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=5242880,RESOLUTION=1280x720,CODECS="avc1.42e00a,mp4a.40.2"
720p/index.m3u8
`)

	return ioutil.WriteFile(filename, m3u8, 0755)
}

func waitForStreamReady(streamurl string) {
	maxretry := 10
	for i := 0; i < maxretry; i++ {
		resp, _ := http.Head(streamurl)
		if resp.StatusCode == 200 {
			return
		}
		log.Infof("waiting for stream %s to become ready...", streamurl)
		time.Sleep(30 * time.Second)
	}
}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0777)
}

func buildCmd(inputURL string, dir string) []string {
	cmd := []string{"-re", "-i", inputURL}

	for _, b := range bitrates {
		args := fmt.Sprintf("-hls_allow_cache 0 -hls_flags append_list -f ssegment -b:v %d -strict -2 -c:v h264 -profile:v main -segment_list_flags live -segment_time 10 -segment_format mpegts -an -segment_list %s/%d/index.m3u8 %s/%d/%%d.ts", b, dir, b, dir, b)
		cmd = append(cmd, strings.Split(args, " ")...)

	}

	return cmd

}
