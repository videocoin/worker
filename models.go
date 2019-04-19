package transcode

import (
	"context"
	"math/big"
	"os/exec"

	manager_v1 "github.com/VideoCoin/cloud-api/manager/v1"
	verifier_v1 "github.com/VideoCoin/cloud-api/verifier/v1"
	"github.com/VideoCoin/cloud-pkg/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/grafov/m3u8"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

// Byte constants
const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

var (
	bitrates = []uint32{
		8 * MB,
	}
)

type (
	// Job used to queue up chunk work
	Job struct {
		Bitrate         uint32
		InputChunkName  string
		OutputChunkName string
		ChunksDir       string
		ContractAddr    string
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
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		pkAddr        common.Address
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
	}
)
