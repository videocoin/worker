package transcoder

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Transcoder struct {
	logger *logrus.Entry
	t      *time.Ticker
}

func NewTranscoder(logger *logrus.Entry) (*Transcoder, error) {
	return &Transcoder{
		logger: logger,
	}, nil
}

func (t *Transcoder) Start() error {
	t.logger.Infof("starting transcoder")
	t.t = time.NewTicker(5 * time.Second)
	for range t.t.C {
		t.logger.Infof("getting task...")
		time.Sleep(time.Second * 2)
		t.logger.Infof("no task")
	}
	return nil
}

func (t *Transcoder) Stop() error {
	t.logger.Infof("stopping transcoder")
	t.t.Stop()
	return nil
}
