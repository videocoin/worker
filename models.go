package transcode

import (
	"context"
	"math/big"

	"github.com/VideoCoin/common/proto"
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
	// CSync struct for handling sync logic

	// Job used to queue up chunk work
	Job struct {
		Bitrate         uint32
		InputChunkName  string
		OutputChunkName string
		ChunksDir       string
		StreamID        *big.Int
		InputID         *big.Int
		OutputID        *big.Int
		Playlist        *m3u8.MediaPlaylist
		Wallet          common.Address
	}

	// Service primary reciever for service
	Service struct {
		cfg           *Config
		ctx           context.Context
		manager       proto.ManagerServiceClient
		verifier      proto.VerifierServiceClient
		streamManager *streamManager.Manager
		bcAuth        *bind.TransactOpts
		bcClient      *ethclient.Client
		pkAddr        common.Address
		log           *logrus.Entry
	}
)
