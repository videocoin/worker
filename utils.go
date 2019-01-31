package transcode

import (
	"context"
	"crypto/rand"
	"io"
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
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(64), nil).Sub(max, big.NewInt(1))
	n, _ := rand.Int(rand.Reader, max)
	return n
}

// ChunkNum strip .ts from input chunk and return as bigInt
func ChunkNum(chunkName string) *big.Int {
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
	if err != nil {
		return 0.0, err
	}

	cleanOut := strings.TrimSpace(string(stdout))

	return strconv.ParseFloat(cleanOut, 64)
}
