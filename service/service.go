package service

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	minersv1 "github.com/videocoin/cloud-api/miners/v1"
	syncerv1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/pinger"
	"github.com/videocoin/transcode/transcoder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Service struct {
	cfg        *Config
	dispatcher dispatcherv1.DispatcherServiceClient
	transcoder *transcoder.Transcoder
	syncer     syncerv1.SyncerServiceClient
	pinger     *pinger.Pinger

	ClientID string
}

func NewService(cfg *Config) (*Service, error) {
	cli, err := ethclient.Dial(cfg.RPCNodeURL)
	if err != nil {
		return nil, err
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, cli)
	if err != nil {
		return nil, err
	}

	dGrpcDialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 10,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		}),
	}

	dispatcherConn, err := grpc.Dial(cfg.DispatcherRPCAddr, dGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(dispatcherConn)

	_, err = dispatcher.Register(
		context.Background(),
		&minersv1.RegistrationRequest{
			ClientID: cfg.ClientID,
			Address:  caller.Addr().String(),
		},
	)
	if err != nil {
		return nil, err
	}

	translogger := cfg.Logger
	trans, err := transcoder.NewTranscoder(
		translogger,
		dispatcher,
		cfg.ClientID,
		cfg.OutputDir,
		caller,
		cfg.SyncerURL,
	)
	if err != nil {
		return nil, err
	}

	plogger := cfg.Logger.WithField("system", "pinger")
	pinger, err := pinger.NewPinger(dispatcher, cfg.ClientID, time.Second*5, cfg.Version, plogger)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:        cfg,
		dispatcher: dispatcher,
		transcoder: trans,
		pinger:     pinger,
		ClientID:   cfg.ClientID,
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.transcoder.Start()
	go s.pinger.Start()
	return nil
}

func (s *Service) Stop() error {
	s.transcoder.Stop()
	s.pinger.Stop()
	return nil
}
