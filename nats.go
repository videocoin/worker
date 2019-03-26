package transcode

import (
	"encoding/json"

	pb "github.com/VideoCoin/common/proto"
	"github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
)

const (
	queueGroup = "transcoder"
)

func connectNats(natsURL, clusterID, clientID string) (stan.Conn, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		return sc, err
	}

	return sc, nil
}

func (s *Service) subscribe(uid string) {
	s.listenForAssignment(uid)
}

func (s *Service) listenForAssignment(uid string) {

	qcb := func(m *stan.Msg) {
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
