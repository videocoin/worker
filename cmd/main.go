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
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
	"github.com/videocoin/cloud-pkg/ethutils"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-pkg/logger"
	bridge "github.com/videocoin/go-bridge/client"
	staking "github.com/videocoin/go-staking"
	vcoauth2 "github.com/videocoin/oauth2/videocoin"
	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/service"
	"golang.org/x/oauth2"
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

func validateAmountOps(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		if cmd.Name() == "fund" || cmd.Name() == "add" || cmd.Name() == "withdraw" {
			return errors.New("requires an amount (tokens value) argument")
		}
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(args[0], 10)
	if !ok {
		return errors.New("amount value must be integer")
	}

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("amount value has to be positive")
	}

	return nil
}

func validateDelegateOps(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "delegate" {
		if len(args) < 1 {
			return errors.New("requires worker address argument")
		}
	} else if len(args) < 2 {
		if cmd.Name() == "add" || cmd.Name() == "withdraw" {
			return errors.New("requires worker address and amount (tokens value) arguments")
		}
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

	var startCmd = &cobra.Command{
		Use:              "start",
		Short:            "start worker",
		TraverseChildren: true,
		PreRunE:          validateFlags,
		Run:              runStartCommand,
	}

	var fundCmd = &cobra.Command{
		Use:              "fund",
		Short:            "fund native account",
		TraverseChildren: true,
		Args:             validateAmountOps,
		PreRunE:          validateFlags,
		Run:              runFundCommand,
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
		Short:            "add stake of specified VideoCoin token amount",
		TraverseChildren: true,
		Args:             validateAmountOps,
		PreRunE:          validateFlags,
		Run:              runStakeAddCommand,
	}

	var withdrawCmd = &cobra.Command{
		Use:              "withdraw",
		Short:            "withdraw stake of specified tokens amount",
		TraverseChildren: true,
		Args:             validateAmountOps,
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
		Short:            "delegate stake of specified VideoCoin token amount to worker",
		TraverseChildren: true,
		Args:             validateDelegateOps,
		PreRunE:          validateFlags,
		Run:              runDelegateAddCommand,
	}

	var dwithdrawCmd = &cobra.Command{
		Use:              "withdraw",
		Short:            "withdraw delegate stake of specified tokens amount from worker",
		TraverseChildren: true,
		Args:             validateDelegateOps,
		PreRunE:          validateFlags,
		Run:              runDelegateWithdrawCommand,
	}

	var withdrawCompleteCmd = &cobra.Command{
		Use:              "complete",
		Short:            "complete pending withdraw after unbonding period finish",
		TraverseChildren: true,
		Args:             nil,
		PreRunE:          validateFlags,
		Run:              runWithdrawCompleteCommand,
	}

	// root command initialize
	rootCmd.PersistentFlags().StringP("loglevel", "l", "INFO", "")
	rootCmd.PersistentFlags().StringP("key", "k", "", "utc key file json content")
	rootCmd.PersistentFlags().StringP("secret", "s", "", "password to decrypt key file")

	// start command initialize
	startCmd.Flags().StringP("client-id", "c", "", "unique client id assigned to worker (required)")

	// add commands and execute
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(fundCmd)
	rootCmd.AddCommand(stakeCmd)
	stakeCmd.AddCommand(withdrawCmd)
	stakeCmd.AddCommand(addCmd)
	rootCmd.AddCommand(delegateCmd)
	delegateCmd.AddCommand(dwithdrawCmd)
	delegateCmd.AddCommand(daddCmd)

	withdrawCmd.AddCommand(withdrawCompleteCmd)
	dwithdrawCmd.AddCommand(withdrawCompleteCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func runStartCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)
	cfg.Logger = log

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

func getClients(cfg *service.Config) (*staking.Client, *bridge.Client, *caller.Caller, error) {
	conn, err := grpcutil.Connect(cfg.DispatcherRPCAddr, cfg.Logger.WithField("system", "dispatcher"))
	if err != nil {
		return nil, nil, nil, err
	}
	dispatcher := dispatcherv1.NewDispatcherServiceClient(conn)

	configReq := new(dispatcherv1.ConfigRequest)
	configResp, err := dispatcher.GetDelegatorConfig(
		context.Background(),
		configReq,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	symphonyTS, err := vcoauth2.JWTAccessTokenSourceFromJSON([]byte(configResp.AccessKey), configResp.RPCNodeURL)
	if err != nil {
		return nil, nil, nil, err
	}

	symphonyCli := oauth2.NewClient(context.Background(), symphonyTS)
	symphonyRPCCli, err := ethrpc.DialHTTPWithClient(configResp.RPCNodeURL, symphonyCli)
	if err != nil {
		return nil, nil, nil, err
	}

	natClient := ethclient.NewClient(symphonyRPCCli)
	ethClient, err := ethclient.Dial(cfg.ETHNodeURL)
	if err != nil {
		return nil, nil, nil, err
	}

	caller, err := caller.NewCaller(cfg.Key, cfg.Secret, natClient)
	if err != nil {
		return nil, nil, nil, err
	}

	stakingClient, err := staking.NewClient(natClient, common.HexToAddress(cfg.StakingManagerAddr))
	if err != nil {
		return nil, nil, nil, err
	}

	bridgeClient, err := bridge.Dial(natClient, ethClient, bridge.Config{
		ProxyAddress:         common.HexToAddress(cfg.ProxyAddr),
		ERC20Address:         common.HexToAddress(cfg.ERC20Addr),
		LocalBridgeAddress:   common.HexToAddress(cfg.LocalBridgeAddr),
		ForeignBridgeAddress: common.HexToAddress(cfg.ForeignBridgeAddr),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return stakingClient, bridgeClient, caller, nil
}

func runFundCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	_, bClient, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	tokenAmount := new(big.Int)
	tokenAmount, ok := tokenAmount.SetString(args[0], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	fAmount, _ := new(big.Float).SetInt(tokenAmount).Float64()
	weiAmount := ethutils.EthToWei(fAmount)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	_, err = bClient.WaitDeposit(ctx, caller.PrivateKey(), common.HexToAddress(cfg.TokenBankAddr), weiAmount)
	if err != nil {
		log.Fatal(err.Error())
	}
	cancel()

	log.Infof("account %s has been successfully funded on native network", caller.Addr().String())
}

func runStakeCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	sClient, _, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	t, err := sClient.GetTranscoder(context.Background(), caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenAmount, _ := ethutils.WeiToEth(t.SelfStake)
	log.Infof("worker is staking %s tokens and state is %v", tokenAmount.String(), t.State)
}

func runStakeAddCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	sClient, bClient, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	tokenAmount := new(big.Int)
	tokenAmount, ok := tokenAmount.SetString(args[0], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	fAmount, _ := new(big.Float).SetInt(tokenAmount).Float64()
	weiAmount := ethutils.EthToWei(fAmount)

	if err := sClient.RegisterTranscoder(context.Background(), caller.PrivateKey(), 0); err != nil && !errors.Is(err, staking.ErrAlreadyRegistered) {
		log.Fatal(err.Error())
	}

	t, err := sClient.GetTranscoder(context.Background(), caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	requiredStakeAmount, err := sClient.GetRequiredSelfStake(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	if t.SelfStake.Uint64() == 0 && requiredStakeAmount.Cmp(weiAmount) > 0 {
		log.Fatal("stake amount is insufficient")
	}

	// TODO: check for delegate stake
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	_, err = bClient.WaitDeposit(ctx, caller.PrivateKey(), common.HexToAddress(cfg.TokenBankAddr), weiAmount)
	if err != nil {
		log.Fatal(err.Error())
	}
	cancel()

	if err := sClient.Delegate(context.Background(), caller.PrivateKey(), caller.Addr(), weiAmount); err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %s tokens has been successfully added", tokenAmount.String())
}

func runStakeWithdrawCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	sClient, _, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenAmount := new(big.Int)
	tokenAmount, ok := tokenAmount.SetString(args[0], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	t, err := sClient.GetTranscoder(context.Background(), caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	fAmount, _ := new(big.Float).SetInt(tokenAmount).Float64()
	weiAmount := ethutils.EthToWei(fAmount)

	if weiAmount.Cmp(t.SelfStake) > 0 {
		log.Fatal(fmt.Errorf("amount to withdraw is bigger than existent stake"))
	}

	if _, err := sClient.RequestWithdrawal(context.Background(), caller.PrivateKey(), caller.Addr(), weiAmount); err != nil {
		log.Fatal(err.Error())
	}

	unbondingTime, _ := sClient.GetUnbondingPeriod(context.Background())
	log.Infof("stake withdraw of amount %s tokens has been successfully submitted and be available to complete after %s unbonding periods", tokenAmount.String(), unbondingTime.String())

	if t.State != staking.StateBonded {
		log.Infof("withdraw has been finished")
	} else {
		log.Infof("finish withdraw with `stake withdraw complete` command")
	}
}

func runDelegateCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	sClient, _, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := common.HexToAddress(args[0])
	stake, err := sClient.GetDelegatorStake(context.Background(), addr, caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenAmount, _ := ethutils.WeiToEth(stake)
	log.Infof("delegator is staking %s tokens", tokenAmount.String())
}

func runDelegateAddCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	sClient, bClient, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := common.HexToAddress(args[0])

	tokenAmount := new(big.Int)
	tokenAmount, ok := tokenAmount.SetString(args[1], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	fAmount, _ := new(big.Float).SetInt(tokenAmount).Float64()
	weiAmount := ethutils.EthToWei(fAmount)

	minAmount, err := sClient.GetMinDelegation(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	if minAmount.Cmp(weiAmount) > 0 {
		log.Fatalf("minimum amount to delegate is %s", minAmount.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	_, err = bClient.WaitDeposit(ctx, caller.PrivateKey(), common.HexToAddress(cfg.TokenBankAddr), weiAmount)
	if err != nil {
		log.Fatal(err.Error())
	}
	cancel()

	if err := sClient.Delegate(context.Background(), caller.PrivateKey(), addr, weiAmount); err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("stake of amount %s tokens has been successfully delegated", tokenAmount.String())
}

func runDelegateWithdrawCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	sClient, _, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	addr := common.HexToAddress(args[0])

	tokenAmount := new(big.Int)
	tokenAmount, ok := tokenAmount.SetString(args[1], 10)
	if !ok {
		log.Fatal(err.Error())
	}

	t, err := sClient.GetTranscoder(context.Background(), addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	if t.State >= staking.StateUnbonded {
		log.Fatal("transcoder is unbonded")
	}

	fAmount, _ := new(big.Float).SetInt(tokenAmount).Float64()
	weiAmount := ethutils.EthToWei(fAmount)

	stake, err := sClient.GetDelegatorStake(context.Background(), addr, caller.Addr())
	if err != nil {
		log.Fatal(err.Error())
	}

	if weiAmount.Cmp(stake) > 0 {
		log.Fatal(fmt.Errorf("amount to withdraw is bigger than existent stake"))
	}

	if _, err := sClient.RequestWithdrawal(context.Background(), caller.PrivateKey(), addr, weiAmount); err != nil {
		log.Fatal(err.Error())
	}

	unbondingTime, _ := sClient.GetUnbondingPeriod(context.Background())
	log.Infof("stake withdraw of amount %s tokens has been successfully submitted and be available to complete after %s unbonding periods", tokenAmount.String(), unbondingTime.String())

	if t.State != staking.StateBonded {
		log.Infof("withdraw has been finished")
	} else {
		log.Infof("finish withdraw with `delegate withdraw complete` command")
	}
}

func runWithdrawCompleteCommand(cmd *cobra.Command, args []string) {
	log := logger.NewLogrusLogger(ServiceName, Version)

	sClient, bClient, caller, err := getClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	winfo, err := sClient.CompleteWithdrawals(context.Background(), caller.PrivateKey())
	if err != nil {
		if errors.Is(err, staking.ErrNoPendingWithdrawals) {
			log.Infof("there are no pending withdrawals")
			return
		}

		log.Fatal(err.Error())
	}

	if winfo.Amount.Cmp(big.NewInt(0)) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		_, err := bClient.WaitWithdraw(ctx, caller.PrivateKey(), common.HexToAddress(cfg.NativeBankAddr), winfo.Amount)
		if err != nil {
			log.Fatal(err.Error())
		}

		fVidValue, _ := ethutils.WeiToEth(winfo.Amount)
		vidValue, _ := fVidValue.Int(nil)

		log.Infof("stake withdraw of amount %s tokens have been successfully completed", vidValue.String())
	}
}
