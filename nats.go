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
	{
		s.ec.BindRecvChan(uid, assignmentCh)
		s.listenForAssignment(uid)
		s.heartBeat(uid)
	}
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
	s.ec.Subscribe(fmt.Sprintf("%s-ping", uid), func(subj, reply, msg string) {
		s.ec.Publish(reply, "pong")
	})
}
