package transcode

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"

	nats "github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
	log "github.com/sirupsen/logrus"
)

// Service base struct for service reciever
type Service struct {
	cfg *Config
	sc  stan.Conn
}

// New initialize and return a new Service object
func New() (*Service, error) {
	cfg := LoadConfig()

	// Generate unique connection name
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	clientID := hex.EncodeToString(b)

	// Use nats as underlying connection with secret token
	nc, err := nats.Connect(cfg.NATsURL, nats.Token(cfg.NATsToken))
	if err != nil {
		return nil, err
	}

	// Wrap stan connection ontop of nats
	sc, err := stan.Connect(cfg.Cluster, clientID, stan.NatsConn(nc))
	if err != nil {
		return nil, err
	}

	return &Service{
		cfg: cfg,
		sc:  sc,
	}, nil

}

// Start creates new service and blocks until stop signal
func Start() {
	var inputURL = flag.String("input_url", "", "input stream url")
	var taskID = flag.String("task_id", "", "id of task")
	s, err := New()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	if taskID == nil || *taskID == "" {
		flag.PrintDefaults()
		return
	}

	if inputURL == nil || *inputURL == "" {
		flag.PrintDefaults()
		return
	}

	s.handleTranscodeTask(*taskID, *inputURL)

	handleExit()
}

func (s *Service) handleTranscodeTask(taskID string, inputURL string) error {

	log.Infof("starting transcode task:\n%+s using input: %s", taskID, inputURL)

	dir := path.Join(s.cfg.OutputDir, taskID)

	if err := prepareDir(dir); err != nil {
		panic(err)
	}

	if err := generatePlaylist(path.Join(dir, "playlist.m3u8")); err != nil {
		panic(err)
	}

	args := buildCmd(inputURL, dir)

	transcode(args)

	return nil
}

func transcode(args []string) error {
	log.Info("starting transcode")
	out, err := exec.Command("ffmpeg", args...).CombinedOutput()
	if err != nil {
		log.Errorf("failed to exec - output: %s", string(out))
		return err
	}
	log.Info("transcode complete")
	return nil
}

func generatePlaylist(filename string) error {
	m3u8 := []byte(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360
360p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=1400000,RESOLUTION=842x480
480p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2800000,RESOLUTION=1280x720
720p.m3u8
`)

	return ioutil.WriteFile(filename, m3u8, 0644)
}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0666)
}

func buildCmd(inputURL string, dir string) []string {

	p360 := fmt.Sprintf("-vf scale=w=640:h=360:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -b:v 800k -maxrate 856k -bufsize 1200k -hls_segment_filename %s/360p_%%03d.ts %s/360p.m3u8", dir, dir)
	p480 := fmt.Sprintf("-vf scale=w=842:h=480:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -b:v 1400k -maxrate 1498k -bufsize 2100k -hls_segment_filename %s/480p_%%03d.ts %s/480p.m3u8", dir, dir)
	p720 := fmt.Sprintf("-vf scale=w=1280:h=720:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -b:v 2800k -maxrate 2996k -bufsize 4200k -hls_segment_filename %s/720p_%%03d.ts %s/720p.m3u8", dir, dir)
	p1080 := fmt.Sprintf("-vf scale=w=1920:h=1080:force_original_aspect_ratio=decrease -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 10 -hls_playlist_type event -b:v 5000k -maxrate 5350k -bufsize 7500k -hls_segment_filename %s/1080p_%%03d.ts %s/1080p.m3u8", dir, dir)

	cmd := []string{"-i", inputURL}
	cmd = append(cmd, strings.Split(p360, " ")...)
	cmd = append(cmd, strings.Split(p480, " ")...)
	cmd = append(cmd, strings.Split(p720, " ")...)
	cmd = append(cmd, strings.Split(p1080, " ")...)

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
