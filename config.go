package transcode

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gitlab.videocoin.io/videocoin/common/models"
)

var cfg models.Transcoder
var once sync.Once

// LoadConfig initialize config
func LoadConfig(loc string) *models.Transcoder {
	switch loc {
	case "local":
		once.Do(func() {
			err := envconfig.Process("", &cfg)
			if err != nil {
				logrus.Fatalf("failed to load config: %s", err.Error())
			}
		})
		break
	// requires PROJECT_ID environment variable
	case "remote":
		once.Do(func() {
			ctx := context.Background()
			client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
			if err != nil {
				logrus.Fatalf("failed to create new client: %s", err)
			}

			key := datastore.NameKey("config", "transcoder", nil)
			err = client.Get(ctx, key, &cfg)
			if err != nil {
				logrus.Fatalf("failed to get namekey: %s", err)
			}
		})

		break

	default:

	}

	return &cfg
}
