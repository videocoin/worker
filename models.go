package transcode

import (
	"context"
	"math/big"
	"os/exec"

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

const (
	managerAPIURL  = "manager:50053/api/v1/"
	verifierAPIURL = "verifier:50054/api/v1/"
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
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
	}
)
