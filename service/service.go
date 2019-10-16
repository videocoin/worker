package service

import (
	"strings"

	"github.com/google/uuid"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	syncerv1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/transcode/blockchain"
	"github.com/videocoin/transcode/transcoder"
	"google.golang.org/grpc"
)

type Service struct {
	cfg        *Config
	dispatcher dispatcherv1.DispatcherServiceClient
	transcoder *transcoder.Transcoder
	syncer     syncerv1.SyncerServiceClient

	ClientID string
}

func NewService(cfg *Config) (*Service, error) {
	dlogger := cfg.Logger.WithField("system", "dispatchercli")
	dGrpcDialOpts := grpcutil.ClientDialOptsWithRetry(dlogger)
	dispatcherConn, err := grpc.Dial(cfg.DispatcherRPCAddr, dGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(dispatcherConn)

	clientID := uuid.New()

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
		clientID.String(),
		cfg.OutputDir,
		bccli,
	)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:        cfg,
		dispatcher: dispatcher,
		transcoder: trans,
		ClientID:   clientID.String(),
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.transcoder.Start()
	// go telegraf.Run(
	// 	s.cfg.Logger.WithField("system", "telegraf"),
	// 	s.ClientID)

	return nil
}

func (s *Service) Stop() error {
	s.transcoder.Stop()
	return nil
}
