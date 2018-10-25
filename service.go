package transcode

import (
	"gitlab.videocoin.io/videocoin/common/mqmux"
	"gitlab.videocoin.io/videocoin/transcode/config"
)

// Service base struct for service reciever
type Service struct {
	mq  *mqmux.WorkerMux
	cfg *config.Config
}

// New initialize and return a new Service object
func New() (*Service, error) {
	cfg := config.Load()

	mqmux, err := mqmux.NewWorkerMux(cfg.MqURI, "transcoder")
	if err != nil {
		return nil, err
	}
	return &Service{
		mq:  mqmux,
		cfg: cfg,
	}, nil
}
