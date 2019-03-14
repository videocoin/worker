package transcode

import (
	"context"
	"math/big"
	"os/exec"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
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
		Playlist        *m3u8.MediaPlaylist
		Wallet          common.Address
		cmd             *exec.Cmd
		stopChan        chan struct{}
	}

	// Service primary reciever for service
	Service struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
)
