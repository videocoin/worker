package transcode

import (
	"fmt"

	transcoder_v1 "github.com/VideoCoin/cloud-api/transcoder/v1"
)

func (s *Service) subscribe(uid string) {
	s.heartBeat(uid)

	s.listenForAssignment(uid)
}

func (s *Service) listenForAssignment(uid string) {
	s.ec.Subscribe(uid, func(a *transcoder_v1.Assignment) {
		err := s.handleTranscode(a.Workorder, a.Profile, uid)
		if err != nil {
			s.log.Errorf("failed to handle transcode: %s", err.Error())
		}
	})
}

func (s *Service) heartBeat(uid string) {
	s.log.Info("registering heartbeat monitor")
	_, err := s.ec.Subscribe(fmt.Sprintf("%s-ping", uid), func(subj, reply, msg string) {
		err := s.ec.Publish(reply, "pong")
		if err != nil {
			s.log.Errorf("failed to publish pong: %s", err.Error())
		}
	})
	if err != nil {
		s.log.Errorf("failed to subscribe to heartBeat: %s", err.Error())
	}
}
