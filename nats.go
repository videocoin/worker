package transcode

import (
	"encoding/json"

	pb "github.com/VideoCoin/common/proto"
	stan "github.com/nats-io/go-nats-streaming"
)

const (
	queueGroup = "transcoder"
)

func (s *Service) subscribe(uid string) {
	s.listenForAssignment(uid)
}

func (s *Service) listenForAssignment(uid string) {
	qcb := func(m *stan.Msg) {
		s.log.Infof("recived nats msg: %d", m.CRC32)

		var assignment = new(pb.Assignment)

		err := json.Unmarshal(m.Data, &assignment)
		if err != nil {
			s.log.Errorf("failed to unmarshal work: %s", err.Error())
		}

		s.handleTranscode(assignment.Workorder, assignment.Profile)
	}

	sub, err := s.sc.QueueSubscribe(uid, queueGroup, qcb)
	if err != nil {
		s.log.Errorf("failed to subscribe: %s", err.Error())
	}

	_ = sub

}
