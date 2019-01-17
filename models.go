package transcode

import (
	"context"

	"github.com/grafov/m3u8"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/videocoin/common/proto"

	"github.com/sirupsen/logrus"
)

type (
	// CSync struct for handling sync logic
	CSync struct {
		cfg *Config
		log *logrus.Entry
	}

	// Job used to queue up chunk work
	Job struct {
		ChunkName string
		Folder    string
		Playlist  *m3u8.MediaPlaylist
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
		ctx     context.Context
		csyc    *CSync
	}
)
