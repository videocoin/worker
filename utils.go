package transcode

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os/exec"
	"strconv"
	"strings"

	"google.golang.org/api/storage/v1"
)

// Upload uploads an object to gcs with publicread acl
func (s *Service) upload(output string, r io.Reader) error {
	svc, err := storage.NewService(context.Background())
	if err != nil {
		return err
	}

	object := &storage.Object{
		Name:         output,
		CacheControl: "public, max-age=0",
	}

	if _, err := svc.Objects.Insert(s.cfg.Bucket, object).Media(r).PredefinedAcl("publicread").Do(); err != nil {
		return err
	}

	return nil
}

// ChunkNum strip .ts from input chunk and return as bigInt
func getChunkNum(chunkName string) *big.Int {
	chunkNum, err := strconv.ParseInt(strings.TrimSuffix(chunkName, ".ts"), 10, 64)
	if err != nil {
		return nil
	}
	return big.NewInt(chunkNum)
}

// Duration use ffmpeg to find chunk duration
func (s *Service) duration(input string) (float64, error) {
	args := []string{"-v", "panic", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", input}
	stdout, err := exec.Command("ffprobe", args...).CombinedOutput()
	cleanOut := strings.TrimSpace(string(stdout))
	if err != nil {
		s.log.Warnf("failed to get duraton: %s from: %s", cleanOut, input)
		return 10.00, nil
	}

	return strconv.ParseFloat(cleanOut, 64)
}

// GeneratePlaylist based on static bitrates
func (s *Service) generatePlaylist(streamHash string, filename string, bitrate uint32) error {
	m3u8 := []byte(fmt.Sprintf(`#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1048576,CODECS="avc1.42e00a,mp4a.40.2"
%d/index.m3u8
`, bitrate))

	err := ioutil.WriteFile(filename, m3u8, 0755)
	if err != nil {
		s.log.Errorf("failed to write file: %s", err.Error())
		return err
	}

	reader := bytes.NewReader(m3u8)

	err = s.upload(fmt.Sprintf("%s/%s", streamHash, "index.m3u8"), reader)
	if err != nil {
		s.log.Errorf("failed to upload: %s", err.Error())
		return err
	}

	return nil
}
