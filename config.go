package transcode

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	BaseStreamURL    string `required:"true" envconfig:"BASE_STREAM_URL" default:"http://ingest.videocoin.io:8080/hls"`
	VerifierHTTPADDR string `required:"true" envconfig:"VERIFIER_HTTP_ADDR" default:"http://verifier:50054/api/v1"`
	BaseStorageURL   string `required:"true" envconfig:"BASE_STORAGE_URL" default:"/tmp"`
	Bucket           string `required:"true" envconfig:"BUCKET" default:"streams.videocoin.network"`
	OutputDir        string `required:"true" envconfig:"OUTPUT_DIR" default:"/opt/mnt/" description:"ffmpeg output dir"`
	ManagerHTTPADDR  string `required:"true" envconfig:"MANAGER_HTTP_ADDR" default:"http://manager:50053/api/v1"`
	BlockchainURL    string `required:"true" envconfig:"BLOCKCHAIN_URL" default:"http://localhost:8545"`
	SMCA             string `required:"true" envconfig:"SMCA" default:"0xEa91ac0B88F84e91e79Caa871d2EB04eF5133721" description:"stream manager contract address"`
	Key              string `required:"true" envconfig:"KEY" default:"keys/transcoder.key"`
	Password         string `required:"true" envconfig:"PASSWORD" default:"transcoder"`
	LogLevel         string `required:"true" envconfig:"LOG_LEVEL" default:"DEBUG"`
	HashKey          string `required:"true" envconfig:"HASH" default:"BEEFFEED"`
	NatsURL          string `required:"true" envconfig:"NATS_URL" default:"nats://localhost:4222"`
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
