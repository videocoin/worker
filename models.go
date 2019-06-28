package transcode

import (
	"context"
	"math/big"
	"os/exec"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/grafov/m3u8"
	"github.com/sirupsen/logrus"
	manager_v1 "github.com/videocoin/cloud-api/manager/v1"
	verifier_v1 "github.com/videocoin/cloud-api/verifier/v1"
	"github.com/videocoin/cloud-pkg/streamManager"
)

// Byte constants
const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

type (
	// Job used to queue up chunk work
	Job struct {
		Id              string
		Bitrate         uint32
		InputChunkName  string
		OutputChunkName string
		ChunksDir       string
		StreamAddress   string
		StreamID        *big.Int
		InputID         *big.Int
		OutputID        *big.Int
		cmd             *exec.Cmd
		stopChan        chan struct{}
		Wallet          common.Address
		Playlist        *m3u8.MediaPlaylist
	}

	// Service primary reciever for service
	Service struct {
		cfg           *Config
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
	}
)
