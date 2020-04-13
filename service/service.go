package service

import (
	"context"
	"time"

	"github.com/videocoin/cloud-pkg/grpcutil"

	"cloud.google.com/go/compute/metadata"
	"github.com/ethereum/go-ethereum/ethclient"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	minersv1 "github.com/videocoin/cloud-api/miners/v1"
	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/capacity"
	"github.com/videocoin/transcode/cryptoinfo"
	"github.com/videocoin/transcode/pinger"
	"github.com/videocoin/transcode/transcoder"
	"golang.org/x/oauth2/google"
	computev1 "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type Service struct {
	cfg        *Config
	dispatcher dispatcherv1.DispatcherServiceClient
	transcoder *transcoder.Transcoder
	pinger     *pinger.Pinger

	ClientID string
}

func NewService(cfg *Config) (*Service, error) {
	conn, err := grpcutil.Connect(cfg.DispatcherRPCAddr, cfg.Logger.WithField("system", "dispatcher"))
	if err != nil {
		return nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(conn)

	configReq := new(dispatcherv1.ConfigRequest)
	configResp, err := dispatcher.GetConfig(
		context.Background(),
		configReq,
	)
	if err != nil {
		return nil, err
	}
	cfg.RPCNodeURL = configResp.RPCNodeURL
	cfg.SyncerURL = configResp.SyncerURL

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
	}

	natClient, err := ethclient.Dial(cfg.RPCNodeURL)
	if err != nil {
		return nil, err
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, natClient, nil)
	if err != nil {
		return nil, err
	}

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

	ci, err := cryptoinfo.NewCryptoInfo(caller, cfg.StakingManagerAddr)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info("performing capacity measurements")
	capacitor := capacity.NewCapacitor(cfg.Internal, trans, cfg.Logger)

	pinger, err := pinger.NewPinger(
		dispatcher,
		capacitor,
		ci,
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
		ClientID:   cfg.ClientID,
	}

	return svc, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		errCh <- s.transcoder.Start()
	}()

	s.pinger.Start()
	err := s.markAsRunningOnGCE()
	if err != nil {
		errCh <- err
	}
}

func (s *Service) Stop() error {
	err := s.transcoder.Stop()
	if err != nil {
		return err
	}
	err = s.pinger.Stop()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) markAsRunningOnGCE() error {
	if metadata.OnGCE() {
		project, err := metadata.ProjectID()
		if err != nil {
			s.cfg.Logger.Error(err)
			return err
		}

		zone, err := metadata.Zone()
		if err != nil {
			s.cfg.Logger.Error(err)
			return err
		}

		instanceID, err := metadata.InstanceID()
		if err != nil {
			s.cfg.Logger.Error(err)
			return err
		}

		gctx := context.Background()
		computeCli, err := google.DefaultClient(gctx, computev1.ComputeScope)
		if err != nil {
			s.cfg.Logger.Error(err)
			return err
		}
		ctx := context.Background()
		computeSvc, err := computev1.NewService(ctx, option.WithHTTPClient(computeCli))
		if err != nil {
			s.cfg.Logger.Error(err)
			return err
		}

		instance, err := computeSvc.Instances.Get(project, zone, instanceID).Do()
		if err != nil {
			s.cfg.Logger.Error(err)
			return err
		}

		fingerprint := instance.Metadata.Fingerprint

		md := &computev1.Metadata{
			Fingerprint: fingerprint,
			Items: []*computev1.MetadataItems{
				{
					Key: "vc-running",
				},
			},
		}
		_, err = computeSvc.Instances.SetMetadata(project, zone, instanceID, md).Context(gctx).Do()
		if err != nil {
			s.cfg.Logger.Error(err)
			return err
		}
	}

	return nil
}
