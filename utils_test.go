package transcode

import (
	"context"
	"io"
	"math/big"
	"reflect"
	"testing"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/stream"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/sirupsen/logrus"
)

func TestService_generatePlaylist(t *testing.T) {
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
		filename string
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
			if err := s.generatePlaylist(tt.args.streamID, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Service.generatePlaylist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_duration(t *testing.T) {
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
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
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
			got, err := s.duration(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.duration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getChunkNum(t *testing.T) {
	type args struct {
		chunkName string
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getChunkNum(tt.args.chunkName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getChunkNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomBigInt(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomBigInt(tt.args.len); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("randomBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_upload(t *testing.T) {
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
		output string
		r      io.Reader
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
			if err := s.upload(tt.args.output, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Service.upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
