package transcode

import (
	"context"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/VideoCoin/stream"
	"github.com/VideoCoin/streamManager"
	"github.com/grafov/m3u8"
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
		stream   *stream.Stream
		bcAuth   *bind.TransactOpts
		bcClient *ethclient.Client
		pkAddr   common.Address
	}
)
