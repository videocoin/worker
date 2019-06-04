package transcode

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"strconv"
	"strings"

	manager_v1 "github.com/VideoCoin/cloud-api/manager/v1"
	transcoder_v1 "github.com/VideoCoin/cloud-api/transcoder/v1"
	verifier_v1 "github.com/VideoCoin/cloud-api/verifier/v1"
	"github.com/tidwall/gjson"

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

// RandomBigInt Generates a random big integer with max 64 bits so we can store as int64
func randomBigInt(len int) *big.Int {
	b := make([]byte, len)
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
#EXT-X-VERSION:4
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1048576,RESOLUTION=640x360,CODECS="avc1.42e00a,mp4a.40.2"
%d/index.m3u8
`, bitrate))

	err := ioutil.WriteFile(filename, m3u8, 0755)
	if err != nil {
		s.log.Errorf("failed to write file: %s", err.Error())
		return err
	}

	reader := bytes.NewReader(m3u8)

	err = s.upload(path.Join(streamHash, "index.m3u8"), reader)
	if err != nil {
		s.log.Errorf("failed to upload: %s", err.Error())
		return err
	}

	return nil
}

func checkBalance(address string) (float64, error) {
	response, err := http.Get(path.Join("balance", address))
	if err != nil {
		return 0.0, err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0.0, err
	}

	return gjson.GetBytes(data, "balance").Float(), nil
}

func (s *Service) updateStreamStatus(streamHash, status string) error {
	response, err := http.Post(s.cfg.ManagerHTTPADDR+path.Join(streamHash, status), "application/json", nil)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create post")
	}

	return nil
}

func (s *Service) verify(verifyRequest *verifier_v1.VerifyRequest) error {
	return postForm(s.cfg.VerifierHTTPADDR+"verify", verifyRequest)
}

func (s *Service) registerTranscoder(transcoder *transcoder_v1.Transcoder) error {
	return postForm(s.cfg.ManagerHTTPADDR+"transcoders", transcoder)
}

func (s *Service) registerChunk(chunkRequest *manager_v1.ChunkCreatedRequest) error {
	return postForm(s.cfg.ManagerHTTPADDR+"chunk_created", chunkRequest)
}

func (s *Service) updateTranscoderStatus(id string, status transcoder_v1.TranscoderStatus) error {
	response, err := http.Post(s.cfg.ManagerHTTPADDR+fmt.Sprintf("%s/%d", id, status), "application/json", nil)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update transcoder status: [%d] %s", response.StatusCode, response.Status)
	}

	return nil
}

func postForm(uri string, item interface{}) error {
	form := url.Values{}

	err := encoder.Encode(item, form)
	if err != nil {
		return err
	}

	client := new(http.Client)

	response, err := client.PostForm(uri, form)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post form: [%d] %s", response.StatusCode, response.Status)
	}

	return nil
}
