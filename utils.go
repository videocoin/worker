package transcode

import (
	"context"
	"crypto/rand"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

func sleep() {
	time.Sleep(1 * time.Second)
}

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

func randomBigInt() *big.Int {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(64), nil).Sub(max, big.NewInt(1))
	n, _ := rand.Int(rand.Reader, max)
	return n
}

func getChunkNum(chunkName string) *big.Int {
	chunkNum, err := strconv.ParseInt(strings.TrimSuffix(chunkName, ".ts"), 10, 64)
	if err != nil {
		return nil
	}
	return big.NewInt(chunkNum)
}
