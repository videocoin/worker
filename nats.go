package transcode

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/VideoCoin/common/proto"
	stan "github.com/nats-io/go-nats-streaming"
)

func connect(clusterID string) (
	*NATs,
	error,
) {

	clientID := fmt.Sprintf("%d", time.Now().Unix())
	con, err := stan.Connect(clusterID, clientID)
	if err != nil {
		return nil, err
	}

	return &NATs{
		sc:  con,
		log: logrus.WithField("name", "nats"),
	}, nil

}

func subscribe(
	clusterID string,
	uid string,
) error {
	n, err := connect(clusterID)
	if err != nil {
		return err
	}
	n.work(uid)

	return nil
}

func (n *NATs) work(uid string) {

	var workOrder = new(pb.WorkOrder)

	_, err := n.sc.Subscribe(uid, func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &workOrder)
		if err != nil {
			n.log.Errorf("failed to unmarshal work order: %s", err.Error())
		}
	})

	if err != nil {
		n.sc.Close()
	}

	var sigChan = make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for range sigChan {
			n.sc.Close()
		}
	}()
}
