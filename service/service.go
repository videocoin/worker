package service

import (
	"context"
	"time"

	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	minersv1 "github.com/videocoin/cloud-api/miners/v1"
	syncerv1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/transcode/blockchain"
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
	bcConfig := &blockchain.Config{
		URL:    cfg.BlockchainURL,
		Key:    cfg.Key,
		Secret: cfg.Secret,
		SMCA:   cfg.SMCA,
	}

	bccli, err := blockchain.Dial(bcConfig)
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
			Address:  bccli.RawKey.Address.String(),
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
		bccli,
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
