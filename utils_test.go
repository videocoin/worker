package transcode

import (
	"context"
	"io"
	"math/big"
	"reflect"
	"testing"

	manager_v1 "github.com/VideoCoin/cloud-api/manager/v1"
	transcoder_v1 "github.com/VideoCoin/cloud-api/transcoder/v1"
	verifier_v1 "github.com/VideoCoin/cloud-api/verifier/v1"
	"github.com/VideoCoin/common/streamManager"
	"github.com/VideoCoin/go-videocoin/accounts/abi/bind"
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
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
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
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
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
				ctx:           tt.fields.ctx,
				bcClient:      tt.fields.bcClient,
				bcAuth:        tt.fields.bcAuth,
				streamManager: tt.fields.streamManager,
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
		ctx           context.Context
		bcClient      *ethclient.Client
		bcAuth        *bind.TransactOpts
		streamManager *streamManager.Manager
	}
	type args struct {
		streamHash string
		filename   string
		bitrate    uint32
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
			if err := s.generatePlaylist(tt.args.streamHash, tt.args.filename, tt.args.bitrate); (err != nil) != tt.wantErr {
				t.Errorf("Service.generatePlaylist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkBalance(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkBalance(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_verify(t *testing.T) {
	type args struct {
		v *verifier_v1.VerifyRequest
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
			if err := verify(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateStreamStatus(t *testing.T) {
	type args struct {
		streamHash string
		status     string
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
			if err := updateStreamStatus(tt.args.streamHash, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("updateStreamStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerTranscoder(t *testing.T) {
	type args struct {
		t *transcoder_v1.Transcoder
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
			if err := registerTranscoder(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("registerTranscoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_chunkCreated(t *testing.T) {
	type args struct {
		c *manager_v1.ChunkCreatedRequest
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
			if err := chunkCreated(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("chunkCreated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateTranscoderStatus(t *testing.T) {
	type args struct {
		id     string
		status transcoder_v1.TranscoderStatus
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
			if err := updateTranscoderStatus(tt.args.id, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("updateTranscoderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
