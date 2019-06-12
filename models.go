package transcode

import (
	"context"
	"math/big"
	"os/exec"

	manager_v1 "github.com/videocoin/cloud-api/manager/v1"
	verifier_v1 "github.com/videocoin/cloud-api/verifier/v1"
	"github.com/videocoin/cloud-pkg/streamManager"
	"github.com/videocoin/go-videocoin/accounts/abi/bind"
	"github.com/videocoin/go-videocoin/common"
	"github.com/videocoin/go-videocoin/ethclient"
	"github.com/grafov/m3u8"
	"github.com/sirupsen/logrus"
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
		Bitrate         uint32
		InputChunkName  string
		OutputChunkName string
		ChunksDir       string
		StreamAddress   string
		StreamHash      string
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
