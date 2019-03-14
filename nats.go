package transcode

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	pb "github.com/VideoCoin/common/proto"
	stan "github.com/nats-io/go-nats-streaming"
)

func connectNats(clusterID string) (
	stan.Conn,
	error,
) {

	clientID := fmt.Sprintf("%d", time.Now().Unix())
	con, err := stan.Connect(clusterID, clientID)
	if err != nil {
		return nil, err
	}

	return con, nil

}

func (s *Service) subscribe() {
	s.work()
}

func (s *Service) work() {

	var workOrder = new(pb.WorkOrder)

	_, err := s.sc.Subscribe(s.cfg.UID, func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &workOrder)
		if err != nil {
			s.log.Errorf("failed to unmarshal work order: %s", err.Error())
		}

		profile, err := s.manager.GetProfile(s.ctx, &pb.GetProfileRequest{ProfileId: workOrder.Profile})
		if err != nil {
			s.log.Debugf("failed to get profile: %s", err.Error())
		}

		s.newStream(workOrder.ContractAddress)

		s.handleTranscodeTask(workOrder, profile)

	})

	if err != nil {
		s.sc.Close()
	}

	var sigChan = make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for range sigChan {
			s.sc.Close()
		}
	}()
}
