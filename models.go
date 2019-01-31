package transcode

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/grafov/m3u8"
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
		2 * MB,
		4 * MB,
		8 * MB,
	}
)

type (
	// CSync struct for handling sync logic

	// Job used to queue up chunk work
	Job struct {
		ChunkName string
		ChunksDir string
		Bitrate   uint32
		Playlist  *m3u8.MediaPlaylist
	}

	// JobQueue simple slice of jobs
	JobQueue struct {
		Jobs []Job
	}

	// Service primary reciever for service
	Service struct {
		cfg      *Config
		ctx      context.Context
		manager  proto.ManagerServiceClient
		sm       *streamManager.Manager
		bcAuth   *bind.TransactOpts
		bcClient *ethclient.Client
		pkAddr   common.Address
		log      *logrus.Entry
	}
)
