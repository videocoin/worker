package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/client"
	"github.com/videocoin/transcode/service"
)

var (
	ServiceName string = "transcoder"
	Version     string = "dev"
)

var cfg *service.Config

func validateFlags(cmd *cobra.Command, args []string) error {
	loglevel, _ := cmd.Flags().GetString("loglevel")
	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		loglevel = logrus.InfoLevel.String()
		level, _ = logrus.ParseLevel(loglevel)
	}

	l := logrus.New()
	l.SetLevel(level)
	l.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano})

	cfg.Logger = l.WithField("version", cfg.Version)

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

	if cmd.Name() == "mine" {
		if !cfg.Internal {
			val, err := cmd.Flags().GetString("client-id")
			if val == "" || err != nil {
				if cfg.ClientID == "" {
					return errors.New("client id has to be specified")
				}
			} else {
				cfg.ClientID = val
			}
		}
	}

	return nil
}

func validateAmount(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires an amount argument (wei value)")
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(args[0], 10)
	if !ok {
		return errors.New("amount value must be integer")
	}

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("amount value has to be positive")
	}

	// TODO per command check min, max, etc

	return nil
}

func main() {
	cfg = service.LoadConfig()
	cfg.Name = ServiceName
	cfg.Version = Version

	var rootCmd = &cobra.Command{
		Use: "",
	}

	var mineCmd = &cobra.Command{
		Use:              "mine",
		Short:            "start miner function",
		TraverseChildren: true,
		PreRunE:          validateFlags,
		Run:              runMineCommand,
	}

	var balanceCmd = &cobra.Command{
		Use:              "balance",
		Short:            "show stake balance",
		TraverseChildren: true,
		PreRunE:          validateFlags,
		Run:              runBalanceCommand,
	}

	var stakeCmd = &cobra.Command{
		Use:              "stake",
		Short:            "stake coins",
		TraverseChildren: true,
		Args:             validateAmount,
		PreRunE:          validateFlags,
		Run:              runStakeCommand,
	}

	var withdrawCmd = &cobra.Command{
		Use:              "withdraw",
		Short:            "withdraw staked coins",
		TraverseChildren: true,
		Args:             validateAmount,
		PreRunE:          validateFlags,
		Run:              runWithdrawCommand,
	}

	// root command initialize
	rootCmd.PersistentFlags().StringP("loglevel", "l", "INFO", "")
	rootCmd.PersistentFlags().StringP("key", "k", "", "utc key file json content")
	rootCmd.PersistentFlags().StringP("secret", "s", "", "password to decrypt key file")

	// mine command initialize
	mineCmd.Flags().StringP("client-id", "c", "", "unique client id assigned to miner (required)")

	// withdraw command initialize
	withdrawCmd.Flags().Int64("amount", 0, "amount to withdraw")

	// add commands and execute
	rootCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(stakeCmd)
	rootCmd.AddCommand(withdrawCmd)
	rootCmd.AddCommand(balanceCmd)

	rootCmd.Execute()  //nolint
}

func runMineCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	closer, err := tracer.NewTracer(ServiceName)
	if err != nil {
		log.Info(err.Error())
	} else {
		defer closer.Close()
	}

	svc, err := service.NewService(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		log.Infof("recieved signal %s", sig)
		exit <- true
	}()

	go log.Error(svc.Start())

	<-exit

	err = svc.Stop()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("stopped")
}

func getTranscoderClient(cfg *service.Config) (*client.TranscoderClient, error) {
	cli, err := ethclient.Dial(cfg.RPCNodeURL)
	if err != nil {
		return nil, err
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, cli)
	if err != nil {
		return nil, err
	}

	tCli, err := client.NewTranscoderClient(cfg.StakingManagerAddr, caller)
	if err != nil {
		return nil, err
	}

	return tCli, nil
}

func runStakeCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	cli, err := getTranscoderClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(args[0], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	err = cli.Register(context.Background(), amount)
	if err != nil {
		if err == client.ErrAlreadyRegistered {
			stake, err := cli.GetStake()
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Infof("transcoder is already staking %d wei", stake.Uint64())
		} else {
			log.Fatal(err.Error())
		}
	} else {
		log.Infof("transcoder successfully registered with stake amount of %d wei", amount)
		return
	}

	log.Infof("adding stake in amount of %d", amount)

	err = cli.SelfStake(context.Background(), amount)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %d has been successfully added", amount)
}

func runWithdrawCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	cli, err := getTranscoderClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	stake, err := cli.GetStake()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("transcoder is staking %d wei", stake.Uint64())

	amount := new(big.Int)
	amount, ok := amount.SetString(args[0], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	if amount.Cmp(stake) > 0 {
		log.Fatal(fmt.Errorf("amount to withdraw is bigger than stake"))
	}

	err = cli.WithdrawalProposal(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cli.WithdrawStake(context.Background(), amount)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %d has been successfully withdrawed", amount)
}

func runBalanceCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	cli, err := getTranscoderClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	stake, err := cli.GetStake()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("transcoder is staking %d wei", stake.Uint64())
}
