package service

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	Name               string        `envconfig:"-"`
	Version            string        `envconfig:"-"`
	Logger             *logrus.Entry `envconfig:"-"`
	RPCNodeURL         string        `envconfig:"-"`
	SyncerURL          string        `envconfig:"-"`
	StakingManagerAddr string        `envconfig:"-"`

	DispatcherRPCAddr string `envconfig:"DISPATCHER_ADDR" default:"d.everest.videocoin.network:5008"`
	OutputDir         string `envconfig:"OUTPUT_DIR" default:"/tmp"`
	HealthAddr        string `envconfig:"HEALTH_ADDR" default:"0.0.0.0:8888"`
	LokiURL           string `envconfig:"LOKI_URL"`

	ClientID string `envconfig:"CLIENT_ID"`
	Key      string `envconfig:"KEY"`
	Secret   string `envconfig:"SECRET"`
	Internal bool   `envconfig:"INTERNAL"`
}

var cfg Config
var once sync.Once

func LoadConfig() *Config {
	once.Do(func() {
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatal(err.Error())
		}

	})
	return &cfg
}
