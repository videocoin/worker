package transcode

import (
	"context"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/streamManager"
	"github.com/grafov/m3u8"
	stan "github.com/nats-io/go-nats-streaming"

	"github.com/sirupsen/logrus"
)

type (
	// CSync struct for handling sync logic
	CSync struct {
		manager proto.ManagerServiceClient
		cfg     *Config
		log     *logrus.Entry
		ctx     context.Context
	}

	// Job used to queue up chunk work
	Job struct {
		ChunkName string
		ChunksDir string
		Playlist  *m3u8.MediaPlaylist
		Bitrate   uint32
	}

	// JobQueue simple slice of jobs
	JobQueue struct {
		Jobs []Job
	}

	// Service primary reciever for service
	Service struct {
		cfg     *Config
		sc      stan.Conn
		manager proto.ManagerServiceClient
		bcAuth  *bind.TransactOpts
		sm      *streamManager.Manager
		ctx     context.Context
		csyc    *CSync
	}
)
