package transcode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.videocoin.io/videocoin/common/mqmux"
	"gitlab.videocoin.io/videocoin/common/proto"
	"gitlab.videocoin.io/videocoin/common/tools"
	"gitlab.videocoin.io/videocoin/transcode/config"
)

// Service base struct for service reciever
type Service struct {
	mq  *mqmux.WorkerMux
	cfg *config.Config
}

// New initialize and return a new Service object
func New() (*Service, error) {
	cfg := config.Load()

	mqmux, err := mqmux.NewWorkerMux(cfg.MqURI, "transcoder")
	if err != nil {
		log.Errorf("failed to make new mux worker: %s", err.Error())
		return nil, err
	}

	s := &Service{
		mq:  mqmux,
		cfg: cfg,
	}

	s.mq.Consumer("transcoder", 1, false, s.handleTranscodeTask)

	return s, nil
}

func (s *Service) handleTranscodeTask(d amqp.Delivery) error {
	task := new(proto.TranscodeTask)
	err := json.Unmarshal(d.Body, task)
	if err != nil {
		return err
	}

	log.Infof("starting transcode task:\n%+v", task)

	randomString := tools.RandomID(24)

	dir := path.Join(s.cfg.OutputDir, randomString)

	if err := prepareDir(dir); err != nil {
		return err
	}

	generatePlaylist(path.Join(dir, "playlist.m3u8"))

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
	#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080
	1080p.m3u8`)

	return ioutil.WriteFile(filename, m3u8, 0644)
}

func prepareDir(dir string) error {
	return os.MkdirAll(dir, 0666)
}

func buildCmd(dir string) []string {
	p360 := fmt.Sprintf("-vf scale=w=640:h=360:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod  -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k -hls_segment_filename %s/360p_%%03d.ts %s/360p.m3u8", dir, dir)
	p480 := fmt.Sprintf("-vf scale=w=842:h=480:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 1400k -maxrate 1498k -bufsize 2100k -b:a 128k -hls_segment_filename %s/480p_%%03d.ts %s/480p.m3u8", dir, dir)
	p720 := fmt.Sprintf("-vf scale=w=1280:h=720:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 2800k -maxrate 2996k -bufsize 4200k -b:a 128k -hls_segment_filename %s/720p_%%03d.ts %s/720p.m3u8", dir, dir)
	p1080 := fmt.Sprintf("-vf scale=w=1920:h=1080:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 5000k -maxrate 5350k -bufsize 7500k -b:a 192k -hls_segment_filename %s/1080p_%%03d.ts %s/1080p.m3u8", dir, dir)

	cmd := strings.Split(p360, " ")
	cmd = append(cmd, strings.Split(p480, " ")...)
	cmd = append(cmd, strings.Split(p720, " ")...)
	cmd = append(cmd, strings.Split(p1080, " ")...)

	return cmd

}
