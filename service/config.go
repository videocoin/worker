package service

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	Name    string        `envconfig:"-"`
	Version string        `envconfig:"-"`
	Logger  *logrus.Entry `envconfig:"-"`

	DispatcherRPCAddr  string `required:"true" envconfig:"DISPATCHER_ADDR" default:"d.dev.videocoin.network:5008"`
	RPCNodeURL         string `required:"true" envconfig:"RPC_NODE_URL" default:"https://dev1:D6msEL93LJT5RaPk@rpc.dev.kili.videocoin.network"`
	SyncerURL          string `required:"true" envconfig:"SYNCER_URL" default:"https://dev.videocoin.network/api/v1/sync"`
	OutputDir          string `required:"true" envconfig:"OUTPUT_DIR" default:"/tmp"`
	StakingManagerAddr string `required:"true" envconfig:"STAKING_MANAGER_ADDR" default:"0xeea159afba3969986f8dae95ab1f8eabe6b3ae93"`

	ClientID string `required:"true" envconfig:"CLIENT_ID" default:""`
	Key      string `required:"true" envconfig:"KEY"`
	Secret   string `required:"true" envconfig:"SECRET"`
}

var cfg Config
var once sync.Once

func LoadConfig() *Config {
	once.Do(func() {
		_ = envconfig.Process("", &cfg)
	})

	return &cfg
}
