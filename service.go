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
	stan "github.com/nats-io/go-nats-streaming"
	log "github.com/sirupsen/logrus"
	pb "github.com/videocoin/common/proto"
	"google.golang.org/grpc"
)

// Service base struct for service reciever
type Service struct {
	cfg     *Config
	sc      stan.Conn
	manager pb.ManagerServiceClient
}

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
		return nil, err
	}

	return &Service{
		cfg:     cfg,
		manager: manager,
	}, nil

}

func (s *Service) reportStatus(userID int32, applicationID string, status string) error {
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
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	task, err := s.manager.GetTask(ctx, &pb.GetTaskRequest{Id: hostname})
	if err != nil {
		panic(err)
	}
	//	makePublic(cfg.Bucket, m3u8)

	task.Status = pb.TranscodeStatusTranscoding.String()

	if _, err := s.manager.UpdateStreamStatus(ctx, &pb.UpdateStreamStatusRequest{
		UserId:        task.UserId,
		ApplicationId: task.ApplicationId,
		Status:        task.Status,
	}); err != nil {
		log.Errorf("failed to report status")
	}

	go s.handleTranscodeTask(task)

	handleExit()
}

func (s *Service) handleTranscodeTask(task *pb.SimpleTranscodeTask) error {

	log.Infof("starting transcode task: %+s using input: %s", task.Id, task.InputUrl)

	dir := path.Join(s.cfg.OutputDir, task.Id)
	m3u8 := path.Join(dir, "index.m3u8")

	if err := prepareDir(dir); err != nil {
		log.Error(err.Error())
		// return err
	}

	if err := generatePlaylist(m3u8); err != nil {
		panic(err)
	}

	args := buildCmd(task.InputUrl, dir)

	transcode(args, task.InputUrl)

	return nil
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
360p.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=3145728,RESOLUTION=842x480,CODECS="avc1.42e00a,mp4a.40.2"
480p.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=5242880,RESOLUTION=1280x720,CODECS="avc1.42e00a,mp4a.40.2"
720p.m3u8
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

	// p360 := fmt.Sprintf("-vf scale=w=640:h=360:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -hls_list_size 60 -hls_wrap 10 -start_number 1 -flags -global_header -b:v 800k -maxrate 856k -bufsize 1200k -hls_segment_filename %s/360p_%%03d.ts %s/360p.m3u8", dir, dir)
	// p480 := fmt.Sprintf("-vf scale=w=842:h=480:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -hls_list_size 60 -hls_wrap 10 -start_number 1 -flags -global_header -b:v 1400k -maxrate 1498k -bufsize 2100k -hls_segment_filename %s/480p_%%03d.ts %s/480p.m3u8", dir, dir)
	// p720 := fmt.Sprintf("-vf scale=w=1280:h=720:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -hls_list_size 60 -hls_wrap 10 -start_number 1 -flags -global_header -b:v 2800k -maxrate 2996k -bufsize 4200k -hls_segment_filename %s/720p_%%03d.ts %s/720p.m3u8", dir, dir)
	// /	p1080 := fmt.Sprintf("-vf scale=w=1920:h=1080:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -b:v 5000k -maxrate 5350k -bufsize 7500k -hls_segment_filename %s/1080p_%%03d.ts %s/1080p.m3u8", dir, dir)

	p360 := fmt.Sprintf("-hls_flags delete_segments+append_list -f segment -vf scale=w=640:h=360:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -pix_fmt yuv420p -crf 20 -segment_list_flags live -segment_time 10 -segment_list_size 3 -segment_format mpegts -an -segment_list %s/360p/index.m3u8 %s/360p/%%03d.ts", dir, dir)

	p480 := fmt.Sprintf("-hls_flags delete_segments+append_list -f segment -vf scale=w=842:h=480:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -pix_fmt yuv420p -crf 20 -segment_list_flags live -segment_time 10 -segment_list_size 3 -segment_format mpegts -an -segment_list %s/480p/index.m3u8 %s/480p/%%03d.ts", dir, dir)

	p720 := fmt.Sprintf("-hls_flags delete_segments+append_list -f segment -vf scale=w=1280:h=720:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -pix_fmt yuv420p -crf 20 -segment_list_flags live -segment_time 10 -segment_list_size 3 -segment_format mpegts -an -segment_list %s/720p/index.m3u8 %s/720p/%%03d.ts", dir, dir)

	cmd := []string{"-i", inputURL}
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
