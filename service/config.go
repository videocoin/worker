package service

import (
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	Name    string        `envconfig:"-"`
	Version string        `envconfig:"-"`
	Logger  *logrus.Entry `envconfig:"-"`
	
	DispatcherRPCAddr string `required:"true" envconfig:"DISPATCHER_ADDR" default:"d.snb.videocoin.network:5008"`
	RPCNodeURL        string `required:"true" envconfig:"BLOCKCHAIN_URL" default:"http://admin:VideoCoinS3cr3t@rpc.dev.videocoin.network"`
	SyncerURL         string `required:"true" envconfig:"SYNCER_URL" default:"https://snb.videocoin.network/api/v1/sync"`

	ClientID          string `required:"true" envconfig:"CLIENT_ID"`
	Secret            string `required:"true" envconfig:"SECRET" default:"transcoder"`

	OutputDir         string `required:"true" envconfig:"OUTPUT_DIR" default:"/tmp" description:"local folder for ts chunks"`
	Key               string `required:"true" envconfig:"KEY"`
	
}
