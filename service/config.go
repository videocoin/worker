package service

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	Name       string        `envconfig:"-"`
	Version    string        `envconfig:"-"`
	Logger     *logrus.Entry `envconfig:"-"`
	RPCNodeURL string        `envconfig:"-"`
	SyncerURL  string        `envconfig:"-"`

	DispatcherRPCAddr  string `required:"true" envconfig:"DISPATCHER_ADDR" default:"d.dev.videocoin.network:5008"`
	OutputDir          string `required:"true" envconfig:"OUTPUT_DIR" default:"/tmp"`
	StakingManagerAddr string `required:"true" envconfig:"STAKING_MANAGER_ADDR" default:"0x8Dcbc47852aC4aB74535D375d525a03Ae28E0296"`

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
