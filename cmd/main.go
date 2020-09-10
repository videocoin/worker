package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/videocoin/cloud-api/rpc"
	"github.com/videocoin/cloud-pkg/logger"
	"github.com/videocoin/worker/service"
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

	if !cfg.Internal {
		log = log.Logger.WithField("version", Version)
		log.Logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano})
	}

	cfg.Logger = log

	svc, err := service.NewService(cfg)
	if err != nil {
		if errors.Is(err, rpc.ErrRpcUnauthenticated) {
			log.Fatal("Authentication failed: invalid client-id (-c / --client-id)")
		}
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
