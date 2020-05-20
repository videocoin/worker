package main

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-pkg/logger"
	staking "github.com/videocoin/go-staking"
	vcoauth2 "github.com/videocoin/oauth2/videocoin"
	"github.com/videocoin/worker/caller"
	"github.com/videocoin/worker/service"
	"golang.org/x/oauth2"
)

var (
	ServiceName string = "worker"
	Version     string = "dev"
)

var cfg *service.Config

func validateFlags(cmd *cobra.Command, args []string) error {
	if !cfg.Internal {
		val, err := cmd.Flags().GetString("key")
		if val == "" || err != nil {
			if cfg.Key == "" {
				return errors.New("key file path has to be specified")
			}
		} else {
			if _, err := os.Stat(val); os.IsNotExist(err) {
				return errors.New("key file does not exist")
			}

			keyjson, err := ioutil.ReadFile(val)
			if err != nil {
				return errors.New("failed to read key file")
			}
			cfg.Key = string(keyjson)
		}
	}

	if !cfg.Internal {
		val, err := cmd.Flags().GetString("secret")
		if err != nil {
			return err
		}
		if val != "" {
			cfg.Secret = val
		}
	}

	if cmd.Name() == "start" {
		if !cfg.Internal {
			val, err := cmd.Flags().GetString("client-id")
			if val == "" || err != nil {
				if cfg.ClientID == "" {
					return errors.New("client id has to be specified")
				}
			} else {
				cfg.ClientID = val
			}
		} else {
			val, err := cmd.Flags().GetString("client-id")
			if err != nil && val != "" {
				cfg.ClientID = val
			}
		}
	}

	return nil
}

func main() {
	cfg = service.LoadConfig()
	cfg.Name = ServiceName
	cfg.Version = Version

	var rootCmd = &cobra.Command{
		Use: "",
	}

	var startCmd = &cobra.Command{
		Use:              "start",
		Short:            "start worker",
		TraverseChildren: true,
		PreRunE:          validateFlags,
		Run:              runStartCommand,
	}

	// root command initialize
	rootCmd.PersistentFlags().StringP("loglevel", "l", "INFO", "")
	rootCmd.PersistentFlags().StringP("key", "k", "", "utc key file json content")
	rootCmd.PersistentFlags().StringP("secret", "s", "", "password to decrypt key file")

	// start command initialize
	startCmd.Flags().StringP("client-id", "c", "", "unique client id assigned to worker (required)")

	// add commands and execute
	rootCmd.AddCommand(startCmd)

	err := rootCmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}

func runStartCommand(cmd *cobra.Command, args []string) {
	var lokiURL *string

	if cfg.Internal {
		lokiURL = pointer.ToString(cfg.LokiURL)
	}

	log := logger.NewLogrusLogger(ServiceName, Version, lokiURL)
	cfg.Logger = log

	if !cfg.Internal {
		stakingCli, caller, err := connect(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}

		t, err := stakingCli.GetTranscoder(context.Background(), caller.Addr())
		if err != nil {
			log.Fatal(err.Error())
		}

		if t.State != staking.StateBonded {
			state := ""
			if t.State == staking.StateBonding {
				state = "BONDING"
			} else if t.State == staking.StateUnbonded {
				state = "UNBONDED"
			} else if t.State == staking.StateUnbonding {
				state = "UNBONDING"
			}

			log.Infof("current state is %s", state)
			log.Fatalf("failed to start, state must be BONDED")
		} else {
			log.Info("current state is BONDED")
		}
	}

	svc, err := service.NewService(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	errCh := make(chan error, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		log.Infof("received signal %s", sig)
		_ = svc.Pause()
		exit <- true
	}()

	go svc.Start(errCh)

	select {
	case <-exit:
		break
	case err := <-errCh:
		if err != nil {
			log.Error(err)
		}
		break
	}

	err = svc.Stop()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("stopped")
}

func connect(cfg *service.Config) (*staking.Client, *caller.Caller, error) {
	conn, err := grpcutil.Connect(cfg.DispatcherRPCAddr, cfg.Logger.WithField("system", "dispatcher"))
	if err != nil {
		return nil, nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(conn)

	configReq := new(dispatcherv1.ConfigRequest)
	configResp, err := dispatcher.GetDelegatorConfig(
		context.Background(),
		configReq,
	)
	if err != nil {
		return nil, nil, err
	}

	symphonyTS, err := vcoauth2.JWTAccessTokenSourceFromJSON([]byte(configResp.AccessKey), configResp.RPCNodeURL)
	if err != nil {
		return nil, nil, err
	}

	symphonyCli := oauth2.NewClient(context.Background(), symphonyTS)
	symphonyRPCCli, err := ethrpc.DialHTTPWithClient(configResp.RPCNodeURL, symphonyCli)
	if err != nil {
		return nil, nil, err
	}

	natClient := ethclient.NewClient(symphonyRPCCli)

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, natClient)
	if err != nil {
		return nil, nil, err
	}

	stakingClient, err := staking.NewClient(natClient, common.HexToAddress(cfg.StakingManagerAddr))
	if err != nil {
		return nil, nil, err
	}

	return stakingClient, caller, nil
}
