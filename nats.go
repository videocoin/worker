package transcode

import (
	"fmt"

	pb "github.com/VideoCoin/common/proto"
)

const (
	queueGroup = "transcoder"
)

var (
	assignmentCh = make(chan *pb.Assignment)
)

func (s *Service) subscribe(uid string) {
	s.heartBeat(uid)
	s.ec.BindRecvChan(uid, assignmentCh)
	s.listenForAssignment(uid)
}

func (s *Service) listenForAssignment(uid string) {
	for {
		select {
		case a := <-assignmentCh:
			s.log.Info("recieved assignment")
			s.handleTranscode(a.Workorder, a.Profile)
		}
	}
}

func (s *Service) heartBeat(uid string) {
	s.log.Info("registering heartbeat monitor")
	_, err := s.ec.Subscribe(fmt.Sprintf("%s-ping", uid), func(subj, reply, msg string) {
		s.ec.Publish(reply, "pong")
	})
	if err != nil {
		s.log.Errorf("failed to subscribe to heartBeat: %s", err.Error())
	}
}
