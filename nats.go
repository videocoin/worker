package transcode

import (
	stan "github.com/nats-io/go-nats-streaming"
)

func (s *Service) sub() {

}

func (s *Service) managerHeartbeat() {
	s.sc.Subscribe(s.cfg.UID, func(m *stan.Msg) {
		//
	})
}
