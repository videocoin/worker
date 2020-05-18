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
	ETHNodeURL string        `envconfig:"-" default:"https://goerli.infura.io/v3/300da799d0c54f9bb612088b100ac6ef"`
	SyncerURL  string        `envconfig:"-"`

	DispatcherRPCAddr string `envconfig:"DISPATCHER_ADDR" default:"d.dev.videocoin.network:5008"`
	OutputDir         string `envconfig:"OUTPUT_DIR" default:"/tmp"`
	HealthAddr        string `envconfig:"HEALTH_ADDR" default:"0.0.0.0:8888"`
	LokiURL           string `envconfig:"LOKI_URL"`

	StakingManagerAddr string `envconfig:"STAKING_MANAGER_ADDR" default:"0x74feC37C1CEe00F2EA987080D27e370d79cb46dd"`
	ProxyAddr          string `envconfig:"PROXY_ADDR" default:"0xc16De466447e348b6Cd1B678d604990e6DB3057C"`
	ERC20Addr          string `envconfig:"ERC20_ADDR" default:"0x7bc345BE17AF288a0CFcFF8E673714635C847aa0"`
	LocalBridgeAddr    string `envconfig:"LOCAL_BRIDGE_ADDR" default:"0xb067b9A2eb0bd087F859F836e0AC23E0691Ca62e"`
	ForeignBridgeAddr  string `envconfig:"FOREIGN_BRIDGE_ADDR" default:"0x3CC38A35E3F93B7C57F44330c9584A48ef98E239"`
	TokenBankAddr      string `envconfig:"TOKEN_BANK_ADDR" default:"0x4d80ad6305b893a329039765134ddd436a87ff08"`
	NativeBankAddr     string `envconfig:"NATIVE_BANK_ADDR" default:"0xb8f52379ff40fe8ca57dc60ff24cea17bce043aa"`

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
