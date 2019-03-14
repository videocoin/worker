package transcode

import (
	"context"
	"math/big"
	"os/exec"
	"reflect"
	"testing"

	"github.com/VideoCoin/common/proto"
	pb "github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/core/types"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/sirupsen/logrus"
)

func TestService_handleChunk(t *testing.T) {
	type fields struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
	type args struct {
		job *Job
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
				cfg:            tt.fields.cfg,
				ctx:            tt.fields.ctx,
				manager:        tt.fields.manager,
				log:            tt.fields.log,
				verifier:       tt.fields.verifier,
				streamManager:  tt.fields.streamManager,
				streamInstance: tt.fields.streamInstance,
				bcAuth:         tt.fields.bcAuth,
				bcClient:       tt.fields.bcClient,
				pkAddr:         tt.fields.pkAddr,
			}
			if err := s.handleChunk(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Service.handleChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_syncDir(t *testing.T) {
	type fields struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
	type args struct {
		stop      chan struct{}
		cmd       *exec.Cmd
		workOrder *pb.WorkOrder
		dir       string
		bitrate   uint32
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
				cfg:            tt.fields.cfg,
				ctx:            tt.fields.ctx,
				manager:        tt.fields.manager,
				log:            tt.fields.log,
				verifier:       tt.fields.verifier,
				streamManager:  tt.fields.streamManager,
				streamInstance: tt.fields.streamInstance,
				bcAuth:         tt.fields.bcAuth,
				bcClient:       tt.fields.bcClient,
				pkAddr:         tt.fields.pkAddr,
			}
			s.syncDir(tt.args.stop, tt.args.cmd, tt.args.workOrder, tt.args.dir, tt.args.bitrate)
		})
	}
}

func TestService_submitProof(t *testing.T) {
	type fields struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
	type args struct {
		bitrate       uint32
		inputChunkID  *big.Int
		outputChunkID *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cfg:            tt.fields.cfg,
				ctx:            tt.fields.ctx,
				manager:        tt.fields.manager,
				log:            tt.fields.log,
				verifier:       tt.fields.verifier,
				streamManager:  tt.fields.streamManager,
				streamInstance: tt.fields.streamInstance,
				bcAuth:         tt.fields.bcAuth,
				bcClient:       tt.fields.bcClient,
				pkAddr:         tt.fields.pkAddr,
			}
			got, err := s.submitProof(tt.args.bitrate, tt.args.inputChunkID, tt.args.outputChunkID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.submitProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.submitProof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_verify(t *testing.T) {
	type fields struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
	type args struct {
		tx        *types.Transaction
		job       *Job
		localFile string
		outputURL string
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
				cfg:            tt.fields.cfg,
				ctx:            tt.fields.ctx,
				manager:        tt.fields.manager,
				log:            tt.fields.log,
				verifier:       tt.fields.verifier,
				streamManager:  tt.fields.streamManager,
				streamInstance: tt.fields.streamInstance,
				bcAuth:         tt.fields.bcAuth,
				bcClient:       tt.fields.bcClient,
				pkAddr:         tt.fields.pkAddr,
			}
			if err := s.verify(tt.args.tx, tt.args.job, tt.args.localFile, tt.args.outputURL); (err != nil) != tt.wantErr {
				t.Errorf("Service.verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_process(t *testing.T) {
	type fields struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
	type args struct {
		jobChan   chan Job
		workOrder *pb.WorkOrder
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
				cfg:            tt.fields.cfg,
				ctx:            tt.fields.ctx,
				manager:        tt.fields.manager,
				log:            tt.fields.log,
				verifier:       tt.fields.verifier,
				streamManager:  tt.fields.streamManager,
				streamInstance: tt.fields.streamInstance,
				bcAuth:         tt.fields.bcAuth,
				bcClient:       tt.fields.bcClient,
				pkAddr:         tt.fields.pkAddr,
			}
			s.process(tt.args.jobChan, tt.args.workOrder)
		})
	}
}

func TestService_updateStatus(t *testing.T) {
	type fields struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
	type args struct {
		streamID int64
		status   string
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
				cfg:            tt.fields.cfg,
				ctx:            tt.fields.ctx,
				manager:        tt.fields.manager,
				log:            tt.fields.log,
				verifier:       tt.fields.verifier,
				streamManager:  tt.fields.streamManager,
				streamInstance: tt.fields.streamInstance,
				bcAuth:         tt.fields.bcAuth,
				bcClient:       tt.fields.bcClient,
				pkAddr:         tt.fields.pkAddr,
			}
			s.updateStatus(tt.args.streamID, tt.args.status)
		})
	}
}

func TestService_chunkCreated(t *testing.T) {
	type fields struct {
		cfg            *Config
		ctx            context.Context
		manager        proto.ManagerServiceClient
		log            *logrus.Entry
		verifier       proto.VerifierServiceClient
		streamManager  *streamManager.Manager
		streamInstance *stream.Stream
		bcAuth         *bind.TransactOpts
		bcClient       *ethclient.Client
		pkAddr         common.Address
	}
	type args struct {
		j *Job
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
				cfg:            tt.fields.cfg,
				ctx:            tt.fields.ctx,
				manager:        tt.fields.manager,
				log:            tt.fields.log,
				verifier:       tt.fields.verifier,
				streamManager:  tt.fields.streamManager,
				streamInstance: tt.fields.streamInstance,
				bcAuth:         tt.fields.bcAuth,
				bcClient:       tt.fields.bcClient,
				pkAddr:         tt.fields.pkAddr,
			}
			if err := s.chunkCreated(tt.args.j); (err != nil) != tt.wantErr {
				t.Errorf("Service.chunkCreated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
