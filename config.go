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
	BaseStorageURL string `required:"true" envconfig:"BASE_STORAGE_URL"`
	Bucket         string `required:"true" envconfig:"FUSE_BUCKET"`
	OutputDir      string `required:"true" envconfig:"OUTPUT_DIR" default:"/opt/mnt/" description:"Mount point for GCSFUSE"`
	NATsURL        string `required:"true" envconfig:"NATS_URL" default:"nats://nats:4222"`
	NATsToken      string `required:"true" envconfig:"NATS_TOKEN"`
	Cluster        string `required:"true" envconfig:"NATS_CLUSTER"`
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
