package transcode

import (
	"context"
	"io"
	"math/big"
	"reflect"
	"testing"

	"github.com/VideoCoin/common/proto"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
	"github.com/VideoCoin/go-videocoin/common"
	"github.com/VideoCoin/go-videocoin/ethclient"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

func TestService_upload(t *testing.T) {
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
			if err := s.upload(tt.args.output, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Service.upload() error = %v, wantErr %v", err, tt.wantErr)
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

func TestService_duration(t *testing.T) {
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

func TestService_generatePlaylist(t *testing.T) {
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
			if err := s.generatePlaylist(tt.args.streamID, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Service.generatePlaylist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
