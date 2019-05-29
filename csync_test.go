package transcode

import (
	"context"
	"math/big"
	"os/exec"
	"reflect"
	"testing"

	workorder_v1 "github.com/VideoCoin/cloud-api/workorder/v1"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/core/types"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

func TestService_syncDir(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
	}
	type args struct {
		stop      chan struct{}
		cmd       *exec.Cmd
		workOrder *workorder_v1.WorkOrder
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
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
			}
			s.syncDir(tt.args.stop, tt.args.cmd, tt.args.workOrder, tt.args.dir, tt.args.bitrate)
		})
	}
}

func TestService_handleChunk(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
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
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
			}
			if err := s.handleChunk(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Service.handleChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_submitProof(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
	}
	type args struct {
		contractAddress string
		bitrate         uint32
		inputChunkID    *big.Int
		outputChunkID   *big.Int
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
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
			}
			got, err := s.submitProof(tt.args.contractAddress, tt.args.bitrate, tt.args.inputChunkID, tt.args.outputChunkID)
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
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
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
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
			}
			if err := s.verify(tt.args.tx, tt.args.job, tt.args.localFile, tt.args.outputURL); (err != nil) != tt.wantErr {
				t.Errorf("Service.verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_process(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
	}
	type args struct {
		jobChan   chan Job
		workOrder *workorder_v1.WorkOrder
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
			}
			s.process(tt.args.jobChan, tt.args.workOrder)
		})
	}
}

func TestService_chunkCreated(t *testing.T) {
	type fields struct {
		cfg           *Config
		ec            *nats.EncodedConn
		nc            *nats.Conn
		log           *logrus.Entry
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
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
				cfg:           tt.fields.cfg,
				ec:            tt.fields.ec,
				nc:            tt.fields.nc,
				log:           tt.fields.log,
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
			}
			if err := s.chunkCreated(tt.args.j); (err != nil) != tt.wantErr {
				t.Errorf("Service.chunkCreated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
