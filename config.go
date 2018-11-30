package transcode

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	LogLevel       string `required:"true" default:"debug" envconfig:"LOG_LEVEL" default:"DEBUG"`
	BaseStreamURL  string `required:"true" envconfig:"BASE_STREAM_URL" default:"http://ingester:8080/hls/"`
	BaseStorageURL string `required:"true" envconfig:"BASE_STORAGE_URL" default:"https://storage.googleapis.com/vc-test-fuse"`
	Bucket         string `required:"true" envconfig:"STREAM_BUCKET" default:"vc-test-fuse"`
	OutputDir      string `required:"true" envconfig:"OUTPUT_DIR" default:"/opt/mnt/" description:"Mount point for GCSFUSE"`
	NATsURL        string `required:"true" envconfig:"NATS_URL" default:"nats://nats:4222"`
	NATsToken      string `required:"true" envconfig:"NATS_TOKEN" default:"76cc1e09e6a5c5026ea6868be99f1cb6"`
	Cluster        string `required:"true" envconfig:"NATS_CLUSTER" default:"videocoin"`
}

var cfg Config
var once sync.Once

// LoadConfig initialize config
func LoadConfig() *Config {
	once.Do(func() {
		err := envconfig.Process("", &cfg)
		if err != nil {
			logrus.Fatalf("failed to load config: %s", err.Error())
		}
	})
	return &cfg
}
