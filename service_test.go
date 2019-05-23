package transcode

import (
	"context"
	"os/exec"
	"reflect"
	"testing"

	manager_v1 "github.com/VideoCoin/cloud-api/manager/v1"
	profiles_v1 "github.com/VideoCoin/cloud-api/profiles/v1"
	verifier_v1 "github.com/VideoCoin/cloud-api/verifier/v1"
	workorder_v1 "github.com/VideoCoin/cloud-api/workorder/v1"
	"github.com/VideoCoin/cloud-pkg/stream"
	"github.com/VideoCoin/cloud-pkg/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

func Test_newService(t *testing.T) {
	tests := []struct {
		name    string
		want    *Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newService()
			if (err != nil) != tt.wantErr {
				t.Errorf("newService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_register(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			s.register(tt.args.uid)
		})
	}
}

func TestService_handleTranscode(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
	}
	type args struct {
		workOrder *workorder_v1.WorkOrder
		profile   *profiles_v1.Profile
		uid       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			if err := s.handleTranscode(tt.args.workOrder, tt.args.profile, tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("Service.handleTranscode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_transcode(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
	}
	type args struct {
		cmd          *exec.Cmd
		stop         chan struct{}
		streamurl    string
		contractAddr string
		streamHash   string
		uid          string
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			s.transcode(tt.args.cmd, tt.args.stop, tt.args.streamurl, tt.args.contractAddr, tt.args.streamHash, tt.args.uid)
		})
	}
}

func TestService_waitForStreamReady(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
	}
	type args struct {
		streamurl string
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			s.waitForStreamReady(tt.args.streamurl)
		})
	}
}

func Test_prepareDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := prepareDir(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("prepareDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buildCmd(t *testing.T) {
	type args struct {
		inputURL string
		dir      string
		profile  *profiles_v1.Profile
	}
	tests := []struct {
		name string
		args args
		want *exec.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildCmd(tt.args.inputURL, tt.args.dir, tt.args.profile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_createStreamInstance(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *stream.Stream
		wantErr bool
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			got, err := s.createStreamInstance(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.createStreamInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.createStreamInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_wait(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
		manager       manager_v1.ManagerServiceClient
		verifier      verifier_v1.VerifierServiceClient
	}
	tests := []struct {
		name   string
		fields fields
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
				manager:       tt.fields.manager,
				verifier:      tt.fields.verifier,
			}
			s.wait()
		})
	}
}
