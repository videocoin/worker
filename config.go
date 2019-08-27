package transcode

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	BaseStreamURL   string `required:"true" envconfig:"BASE_STREAM_URL" default:"http://ingest.videocoin.io:8080/hls"`
	VerifierRPCADDR string `required:"true" envconfig:"VERIFIER_RPC_ADDR" default:"verifier:50055"`
	BaseStorageURL  string `required:"true" envconfig:"BASE_STORAGE_URL"`
	Bucket          string `required:"true" envconfig:"BUCKET"`
	OutputDir       string `required:"true" envconfig:"OUTPUT_DIR" default:"/opt/mnt/" description:"local folder for ts chunks"`
	ManagerRPCADDR  string `required:"true" envconfig:"MANAGER_RPC_ADDR" default:"manager:50051"`
	BlockchainURL   string `required:"true" envconfig:"BLOCKCHAIN_URL" default:"http://localhost:8545"`
	SMCA            string `required:"true" envconfig:"SMCA" default:"0xEa91ac0B88F84e91e79Caa871d2EB04eF5133721" description:"stream manager contract address"`
	Key             string `required:"true" envconfig:"KEY"`
	Secret          string `required:"true" envconfig:"SECRET" default:"transcoder"`
	LogLevel        string `required:"true" envconfig:"LOG_LEVEL" default:"DEBUG"`
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
