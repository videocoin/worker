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

	DispatcherRPCAddr string `required:"true" envconfig:"DISPATCHER_ADDR" default:"d.dev.videocoin.network:5008"`
	OutputDir         string `required:"true" envconfig:"OUTPUT_DIR" default:"/tmp"`

	StakingManagerAddr string `required:"true" envconfig:"STAKING_MANAGER_ADDR" default:"0x8Dcbc47852aC4aB74535D375d525a03Ae28E0296"`
	ProxyAddr          string `required:"true" envconfig:"PROXY_ADDR" default:"0x68896bAEcc5186284d5985903e86158F8e803Ca9"`
	ERC20Addr          string `required:"true" envconfig:"ERC20_ADDR" default:"0x7bc345BE17AF288a0CFcFF8E673714635C847aa0"`
	LocalBridgeAddr    string `required:"true" envconfig:"LOCAL_BRIDGE_ADDR" default:"0xfD8c66B99919F291AefE756DA506446FB227f17D"`
	ForeignBridgeAddr  string `required:"true" envconfig:"FOREIGN_BRIDGE_ADDR" default:"0x06C8031ACa8B5d91Ae52F297511Ca4809aB55e29"`
	TokenBankAddr      string `required:"true" envconfig:"TOKEN_BANK_ADDR" default:"0xf90d1852a344a67ca8dabfa1d307b483d89856c4"`
	NativeBankAddr     string `required:"true" envconfig:"NATIVE_BANK_ADDR" default:"0xa2A92CeB62B447A6223fB4f3D2833eE9a40ED9B7"`

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
