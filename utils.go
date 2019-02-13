package transcode

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

// Upload uploads an object to gcs with publicread acl
func (s *Service) Upload(output string, r io.Reader) error {
	client, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}

	svc, err := storage.New(client)
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

// RandomBigInt Generates a random big integer with max 64 bits so we can store as int64
func RandomBigInt() *big.Int {
	b := make([]byte, 8)
	rand.Read(b)
	n := big.NewInt(0)
	n = n.SetBytes(b)

	return n
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
func (s *Service) Duration(input string) (float64, error) {
	s.log.Infof("using input %s", input)
	args := []string{"-v", "panic", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", input}
	stdout, err := exec.Command("ffprobe", args...).CombinedOutput()
	cleanOut := strings.TrimSpace(string(stdout))
	if err != nil {
		s.log.Warnf("failed to get duraton: %s", cleanOut)
		return 10.00, nil
	}

	return strconv.ParseFloat(cleanOut, 64)
}

// GeneratePlaylist based on static bitrates
func (s *Service) GeneratePlaylist(streamID int64, filename string) error {
	m3u8 := []byte(fmt.Sprintf(`#EXTM3U
#EXT-X-VERSION:6
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1048576,RESOLUTION=640x360,CODECS="avc1.42e00a,mp4a.40.2"
%d/index.m3u8
`, bitrates[0]))

	err := ioutil.WriteFile(filename, m3u8, 0755)
	if err != nil {
		s.log.Errorf("failed to write file: %s", err.Error())
		return err
	}

	reader := bytes.NewReader(m3u8)

	err = s.Upload(fmt.Sprintf("%d/%s", streamID, "index.m3u8"), reader)
	if err != nil {
		s.log.Errorf("failed to upload: %s", err.Error())
		return err
	}

	return nil
}
