package transcode

import (
	"context"
	"testing"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

func TestService_subscribe(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		pkAddr        common.Address
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       proto.ManagerServiceClient
		verifier      proto.VerifierServiceClient
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				pkAddr:        tt.fields.pkAddr,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			s.subscribe(tt.args.uid)
		})
	}
}

func TestService_listenForAssignment(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		pkAddr        common.Address
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       proto.ManagerServiceClient
		verifier      proto.VerifierServiceClient
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				pkAddr:        tt.fields.pkAddr,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			s.listenForAssignment(tt.args.uid)
		})
	}
}

func TestService_heartBeat(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		pkAddr        common.Address
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       proto.ManagerServiceClient
		verifier      proto.VerifierServiceClient
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				pkAddr:        tt.fields.pkAddr,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			s.heartBeat(tt.args.uid)
		})
	}
}
