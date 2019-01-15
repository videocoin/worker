package transcode

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	pb "github.com/videocoin/common/proto"
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
	}, nil

}

func (s *Service) reportStatus(userID string, applicationID string, status string) error {
	ctx := context.Background()
	_, err := s.manager.UpdateStreamStatus(ctx, &pb.UpdateStreamStatusRequest{
		UserId:        userID,
		ApplicationId: applicationID,
		Status:        status,
	})

	if err != nil {
		return err
	}

	return nil
}

// Start creates new service and blocks until stop signal
func Start() {
	s, err := New()
	if err != nil {
		panic(err)
	}

	task, err := s.manager.GetJob(s.ctx, &pb.GetJobRequest{})
	if err != nil {
		panic(err)
	}

	task.Status = pb.WorkOrderStatusTranscoding.String()

	if _, err := s.manager.UpdateStreamStatus(s.ctx, &pb.UpdateStreamStatusRequest{
		UserId:        task.UserId,
		ApplicationId: task.ApplicationId,
		Status:        task.Status,
	}); err != nil {
		log.Errorf("failed to report status")
	}

	go s.handleTranscodeTask(task)

	handleExit()
}

func (s *Service) handleTranscodeTask(workOrder *pb.WorkOrder) error {

	log.Infof("starting transcode task: %d using input: %s", workOrder.Id, workOrder.InputUrl)

	dir := path.Join(s.cfg.OutputDir, workOrder.StreamHash)
	m3u8 := path.Join(dir, "index.m3u8")

	if err := prepareDir(dir); err != nil {
		log.Error(err.Error())
	}

	log.Info("monitoring chunks")
	go s.monitorChunks(path.Join(dir, "360p"), workOrder)

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
		UserId:        task.UserId,
		ApplicationId: task.ApplicationId,
		Status:        task.Status,
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

	return ioutil.WriteFile(filename, m3u8, 0644)
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

func makePublic(bucket string, object string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf("failed to get storage client: %s", err.Error())
		return
	}
	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Errorf("failed to make object public: %s", err.Error())
	}

}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0666)
}

func buildCmd(inputURL string, dir string) []string {
	p360 := fmt.Sprintf("-hls_allow_cache 0 -hls_flags append_list -f ssegment -vf scale=640:-2:force_original_aspect_ratio=decrease -strict -2 -c:v h264 -profile:v main -pix_fmt yuv420p -crf 20 -segment_list_flags live -segment_time 10 -segment_format mpegts -an -segment_list %s/360p/index.m3u8 %s/360p/%%d.ts", dir, dir)

	p480 := fmt.Sprintf("-hls_allow_cache 0 -hls_flags append_list -f ssegment -vf scale=842:-2:force_original_aspect_ratio=decrease -strict -2 -c:v h264 -profile:v main -pix_fmt yuv420p -crf 20 -segment_list_flags live -segment_time 10 -segment_format mpegts -an -segment_list %s/480p/index.m3u8 %s/480p/%%d.ts", dir, dir)

	p720 := fmt.Sprintf("-hls_allow_cache 0 -hls_flags append_list -f ssegment -vf scale=1280:-2:force_original_aspect_ratio=decrease -strict -2 -c:v h264 -profile:v main -pix_fmt yuv420p -crf 20 -segment_list_flags live -segment_time 10 -segment_format mpegts -an -segment_list %s/720p/index.m3u8 %s/720p/%%d.ts", dir, dir)

	//p1080 := fmt.Sprintf("-hls_allow_cache 0 -hls_flags append_list -f ssegment -vf scale=1920:-2:force_original_aspect_ratio=decrease -strict -2 -c:v h264 -profile:v main -pix_fmt yuv420p -crf 20 -segment_list_flags live -segment_time 10 -segment_format mpegts -an -segment_list %s/1080p/index.m3u8 %s/1080p/%%03d.ts", dir, dir)

	cmd := []string{"-re", "-i", inputURL}
	cmd = append(cmd, strings.Split(p360, " ")...)
	cmd = append(cmd, strings.Split(p480, " ")...)
	cmd = append(cmd, strings.Split(p720, " ")...)
	//cmd = append(cmd, strings.Split(p1080, " ")...)

	return cmd

}

var (
	done  = make(chan bool)
	errCh = make(chan error)
)

func handleExit() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	sig := signals
	log.Infof("recieved os signal: %v", <-sig)
}
