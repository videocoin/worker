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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-pkg/tracer"
	staking "github.com/videocoin/go-staking"
	"github.com/videocoin/transcode/caller"
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

func validateStakeOps(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires an amount (wei value) argument")
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

func validateDelegateOps(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "delegate" {
		if len(args) < 1 {
			return errors.New("requires worker address argument")
		}
	} else if len(args) < 2 {
		return errors.New("requires worker address and amount (wei value) arguments")
	}

	addr := args[0]
	if !common.IsHexAddress(addr) {
		return errors.New("worker address value is incorrect")
	}

	if cmd.Name() != "delegate" {
		amount := new(big.Int)
		amount, ok := amount.SetString(args[1], 10)
		if !ok {
			return errors.New("amount value must be integer")
		}

		if amount.Cmp(big.NewInt(0)) <= 0 {
			return errors.New("amount value has to be positive")
		}
	}

	return nil
}

func main() {
	cfg = service.LoadConfig()
	cfg.Name = ServiceName
	cfg.Version = Version

	log := logrus.WithFields(logrus.Fields{
		"service": ServiceName,
		"version": Version,
	})

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

	var stakeCmd = &cobra.Command{
		Use:              "stake",
		Short:            "stake operations",
		TraverseChildren: true,
		PreRunE:          validateFlags,
		Run:              runStakeCommand,
	}

	var addCmd = &cobra.Command{
		Use:              "add",
		Short:            "add stake of specified amount",
		TraverseChildren: true,
		Args:             validateStakeOps,
		PreRunE:          validateFlags,
		Run:              runStakeAddCommand,
	}

	var withdrawCmd = &cobra.Command{
		Use:              "withdraw",
		Short:            "withdraw stake of specified amount",
		TraverseChildren: true,
		Args:             validateStakeOps,
		PreRunE:          validateFlags,
		Run:              runStakeWithdrawCommand,
	}

	var delegateCmd = &cobra.Command{
		Use:              "delegate",
		Short:            "delegate operations",
		TraverseChildren: true,
		Args:             validateDelegateOps,
		PreRunE:          validateFlags,
		Run:              runDelegateCommand,
	}

	var daddCmd = &cobra.Command{
		Use:              "add",
		Short:            "delegate stake of specified amount to worker",
		TraverseChildren: true,
		Args:             validateDelegateOps,
		PreRunE:          validateFlags,
		Run:              runDelegateAddCommand,
	}

	var dwithdrawCmd = &cobra.Command{
		Use:              "withdraw",
		Short:            "withdraw delegate stake of specified amount from worker",
		TraverseChildren: true,
		Args:             validateDelegateOps,
		PreRunE:          validateFlags,
		Run:              runDelegateWithdrawCommand,
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
	stakeCmd.AddCommand(withdrawCmd)
	stakeCmd.AddCommand(addCmd)
	rootCmd.AddCommand(delegateCmd)
	delegateCmd.AddCommand(dwithdrawCmd)
	delegateCmd.AddCommand(daddCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
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
	errCh := make(chan error, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		log.Infof("recieved signal %s", sig)
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

func getStakingClient(cfg *service.Config) (*staking.Client, error) {
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

	ethClient, err := ethclient.Dial(cfg.RPCNodeURL)
	if err != nil {
		return nil, err
	}

	return staking.NewClient(ethClient, common.HexToAddress(cfg.StakingManagerAddr))
}

func runStakeCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	client, err := getStakingClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	stake, err := client.GetTranscoderStake(context.Background(), caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("worker is staking %s wei", stake.String())
}

func runStakeAddCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	client, err := getStakingClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(args[0], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	stake, err := client.GetTranscoderStake(context.Background(), caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	requiredStake, err := client.GetRequiredSelfStake(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	if stake.Uint64() == 0 && requiredStake.Cmp(amount) > 0 {
		log.Fatal("stake amount is insufficient")
	}

	if err := client.RegisterTranscoder(context.Background(), caller.PrivateKey(), 0); err != nil && !errors.Is(err, staking.ErrAlreadyRegistered) {
		log.Fatal(err.Error())
	}

	if err := client.Delegate(context.Background(), caller.PrivateKey(), caller.Addr(), amount); err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %s wei has been successfully added", amount.String())
}

func runStakeWithdrawCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	client, err := getStakingClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(args[0], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	stake, err := client.GetTranscoderStake(context.Background(), caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	if amount.Cmp(stake) > 0 {
		log.Fatal(fmt.Errorf("amount to withdraw is bigger than existent stake"))
	}

	if _, err := client.RequestWithdrawal(context.Background(), caller.PrivateKey(), caller.Addr(), amount); err != nil {
		log.Fatal(err.Error())
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if _, err := client.WaitWithdrawalsCompleted(timeoutCtx, caller.PrivateKey()); err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %s has been successfully withdrawn", amount.String())
}

func runDelegateCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	client, err := getStakingClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := common.HexToAddress(args[0])
	stake, err := client.GetDelegatorStake(context.Background(), addr, caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("delegator is staking %s wei", stake.String())
}

func runDelegateAddCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	client, err := getStakingClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := common.HexToAddress(args[0])

	amount := new(big.Int)
	amount, ok := amount.SetString(args[1], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	if err := client.Delegate(context.Background(), caller.PrivateKey(), addr, amount); err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %s wei has been successfully delegated", amount.String())
}

func runDelegateWithdrawCommand(cmd *cobra.Command, args []string) {
	log := cfg.Logger
	client, err := getStakingClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := common.HexToAddress(args[0])

	amount := new(big.Int)
	amount, ok := amount.SetString(args[1], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	stake, err := client.GetDelegatorStake(context.Background(), addr, caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	if amount.Cmp(stake) > 0 {
		log.Fatal(fmt.Errorf("amount to withdraw is bigger than existent stake"))
	}

	if _, err := client.RequestWithdrawal(context.Background(), caller.PrivateKey(), addr, amount); err != nil {
		log.Fatal(err.Error())
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if _, err := client.WaitWithdrawalsCompleted(timeoutCtx, caller.PrivateKey()); err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %s has been successfully withdrawn", amount.String())
}
