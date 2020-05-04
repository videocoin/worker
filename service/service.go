package service

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	minersv1 "github.com/videocoin/cloud-api/miners/v1"
	"github.com/videocoin/cloud-api/rpc"
	vcoauth2 "github.com/videocoin/oauth2/videocoin"
	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/capacity"
	"github.com/videocoin/transcode/health"
	"github.com/videocoin/transcode/pinger"
	"github.com/videocoin/transcode/transcoder"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Service struct {
	cfg        *Config
	dispatcher dispatcherv1.DispatcherServiceClient
	transcoder *transcoder.Transcoder
	pinger     *pinger.Pinger
	health     *health.Health
}

func NewService(cfg *Config) (*Service, error) {
	grpcDialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpclogrus.UnaryClientInterceptor(cfg.Logger)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 10,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		}),
		grpc.WithPerRPCCredentials(rpc.TokenAuth{Token: cfg.ClientID}),
	}

	conn, err := grpc.Dial(cfg.DispatcherRPCAddr, grpcDialOpts...)
	if err != nil {
		return nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(conn)

	if cfg.Internal {
		internalConfigReq := &dispatcherv1.InternalConfigRequest{}

		cfg.Logger.Info("getting internal config")

		internalConfigResp, err := dispatcher.GetInternalConfig(
			context.Background(),
			internalConfigReq,
		)
		if err != nil {
			return nil, err
		}

		cfg.ClientID = internalConfigResp.ClientId
		cfg.Key = internalConfigResp.Key
		cfg.Secret = internalConfigResp.Secret

		grpcDialOpts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(grpclogrus.UnaryClientInterceptor(cfg.Logger)),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                time.Second * 10,
				Timeout:             time.Second * 10,
				PermitWithoutStream: true,
			}),
			grpc.WithPerRPCCredentials(rpc.TokenAuth{Token: cfg.ClientID}),
		}
		conn, err := grpc.Dial(cfg.DispatcherRPCAddr, grpcDialOpts...)
		if err != nil {
			return nil, err
		}
		dispatcher = dispatcherv1.NewDispatcherServiceClient(conn)
	}

	configReq := new(dispatcherv1.ConfigRequest)
	configResp, err := dispatcher.GetConfig(context.Background(), configReq)
	if err != nil {
		return nil, err
	}
	cfg.RPCNodeURL = configResp.RPCNodeURL
	cfg.SyncerURL = configResp.SyncerURL

	symphonyTS, err := vcoauth2.JWTAccessTokenSourceFromJSON([]byte(configResp.AccessKey), cfg.RPCNodeURL)
	if err != nil {
		return nil, err
	}

	symphonyCli := oauth2.NewClient(oauth2.NoContext, symphonyTS)
	symphonyRPCCli, err := ethrpc.DialHTTPWithClient(cfg.RPCNodeURL, symphonyCli)
	if err != nil {
		return nil, err
	}

	natClient := ethclient.NewClient(symphonyRPCCli)
	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, natClient)
	if err != nil {
		return nil, err
	}

	cfg.Logger = cfg.Logger.WithField("cid", cfg.ClientID)
	cfg.Logger.Info("registering")

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

	var capacitor *capacity.Capacitor

	if !cfg.Internal {
		cfg.Logger.Info("performing capacity measurements")
		capacitor = capacity.NewCapacitor(cfg.Internal, trans, cfg.Logger)
	} else {
		capacitor = nil
	}

	pinger, err := pinger.NewPinger(
		dispatcher,
		capacitor,
		cfg.ClientID,
		time.Second*5,
		cfg.Version,
		cfg.Logger.WithField("system", "pinger"))
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:        cfg,
		dispatcher: dispatcher,
		transcoder: trans,
		pinger:     pinger,
	}

	if cfg.Internal {
		svc.health, err = health.NewHealth(cfg.HealthAddr)
		if err != nil {
			return nil, err
		}
	}

	return svc, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		errCh <- s.transcoder.Start()
	}()

	go func() {
		if s.health != nil {
			errCh <- s.health.Start()
		}
	}()

	s.pinger.Start()
}

func (s *Service) Stop() error {
	err := s.transcoder.Stop()
	if err != nil {
		return err
	}

	if s.health != nil {
		err := s.health.Stop()
		if err != nil {
			return err
		}
	}

	err = s.pinger.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Pause() error {
	return s.transcoder.Pause()
}
