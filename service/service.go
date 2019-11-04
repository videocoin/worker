package service

import (
	"context"
	"strings"
	"time"

	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	syncerv1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/transcode/blockchain"
	"github.com/videocoin/transcode/pinger"
	"github.com/videocoin/transcode/transcoder"
	"google.golang.org/grpc"
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
	dlogger := cfg.Logger.WithField("system", "dispatchercli")
	dGrpcDialOpts := grpcutil.DefaultClientDialOpts(dlogger)
	dispatcherConn, err := grpc.Dial(cfg.DispatcherRPCAddr, dGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(dispatcherConn)

	_, err = dispatcher.Register(
		context.Background(),
		&dispatcherv1.RegistrationRequest{ClientID: cfg.ClientID},
	)
	if err != nil {
		return nil, err
	}

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

	translogger := cfg.Logger.
		WithField("system", "transcoder").
		WithField("address", strings.ToLower(bccli.RawKey.Address.String()))

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
