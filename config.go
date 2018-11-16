package transcode

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	LogLevel       string `required:"true" default:"debug" envconfig:"LOG_LEVEL" default:"DEBUG"`
	BaseStreamURL  string `required:"true" envconfig:"BASE_STREAM_URL" default:"http://ingester:1935/hls/"`
	BaseStorageURL string `required:"true" envconfig:"BASE_STORAGE_URL" default:"https://storage.googleapis.com/videocoin-test-streams"`
	OutputDir      string `required:"true" envconfig:"OUTPUT_DIR" default:"/opt/mnt/" description:"Mount point for GCSFUSE"`
	NATsURL        string `required:"true" envconfig:"NATS_URL" default:"nats://nats:4222"`
	NATsToken      string `required:"true" envconfig:"NATS_TOKEN" default:"76cc1e09e6a5c5026ea6868be99f1cb6"`
	Cluster        string `required:"true" envconfig:"NATS_CLUSTER" default:"videocoin"`
}

var cfg Config
var once sync.Once

// LoadConfig initialize config
func LoadConfig(loc string) *Config {
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
