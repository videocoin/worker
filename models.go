package transcode

import (
	"context"

	"github.com/grafov/m3u8"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/videocoin/common/proto"

	"github.com/sirupsen/logrus"
)

// CSync Chunk Sync holds config and logger
type CSync struct {
	cfg *Config
	log *logrus.Entry
}

// Job used for makign a delay buffer for chunk verification
type Job struct {
	ChunkName string
	Folder    string
	Playlist  *m3u8.MediaPlaylist
}

// JobQueue simple slice of jobs
type JobQueue struct {
	Jobs []Job
}

// Service base struct for service reciever
type Service struct {
	cfg     *Config
	sc      stan.Conn
	manager proto.ManagerServiceClient
	ctx     context.Context
}
