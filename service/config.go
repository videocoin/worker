package service

import (
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	Name    string        `envconfig:"-"`
	Version string        `envconfig:"-"`
	Logger  *logrus.Entry `envconfig:"-"`

	DispatcherRPCAddr string `required:"true" envconfig:"DISPATCHER_ADDR" default:"127.0.0.1:5008"`
	OutputDir         string `required:"true" envconfig:"OUTPUT_DIR" default:"/tmp" description:"local folder for ts chunks"`

	// BaseStreamURL   string `required:"true" envconfig:"BASE_STREAM_URL"`
	// VerifierRPCADDR string `required:"true" envconfig:"VERIFIER_RPC_ADDR" default:"127.0.0.1:50055"`
	// BaseStorageURL  string `required:"true" envconfig:"BASE_STORAGE_URL"`
	// Bucket          string `required:"true" envconfig:"BUCKET"`
	// BlockchainURL   string `required:"true" envconfig:"BLOCKCHAIN_URL" default:"http://localhost:8545"`
	// SMCA            string `required:"true" envconfig:"SMCA" default:"0xEa91ac0B88F84e91e79Caa871d2EB04eF5133721" description:"stream manager contract address"`
	// Key             string `required:"true" envconfig:"KEY"`
	// Secret          string `required:"true" envconfig:"SECRET" default:"transcoder"`
}
