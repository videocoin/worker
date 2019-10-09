package service

import (
	"github.com/google/uuid"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/transcode/blockchain"
	"github.com/videocoin/transcode/telegraf"
	"github.com/videocoin/transcode/transcoder"
	"google.golang.org/grpc"
)

type Service struct {
	cfg        *Config
	dispatcher dispatcherv1.DispatcherServiceClient
	transcoder *transcoder.Transcoder

	MachineID string
}

func NewService(cfg *Config) (*Service, error) {
	dlogger := cfg.Logger.WithField("system", "dispatchercli")
	dGrpcDialOpts := grpcutil.ClientDialOptsWithRetry(dlogger)
	dispatcherConn, err := grpc.Dial(cfg.DispatcherRPCAddr, dGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(dispatcherConn)

	machineID := uuid.New()

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

	trans, err := transcoder.NewTranscoder(
		cfg.Logger.WithField("system", "transcoder"),
		dispatcher,
		machineID.String(),
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
		MachineID:  machineID.String(),
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.transcoder.Start()
	go telegraf.Run(
		s.cfg.Logger.WithField("system", "telegraf"),
		s.MachineID)

	return nil
}

func (s *Service) Stop() error {
	s.transcoder.Stop()
	return nil
}
