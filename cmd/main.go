package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	logrussentry "github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	val, err = cmd.Flags().GetString("secret")
	if val != "" {
		cfg.Secret = val
	}

	if cmd.Name() == "mine" {
		val, err = cmd.Flags().GetString("client-id")
		if val == "" || err != nil {
			if cfg.ClientID == "" {
				return errors.New("client id has to be specified")
			}
		} else {
			cfg.ClientID = val
		}
	}

	return nil
}

func validateAmount(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires an amount argument (wei value)")
	}
	if _, err := strconv.Atoi(args[0]); err != nil {
		return errors.New("amount value must be integer")
	}

	// TODO per command check min, max, etc

	return nil
}

func main() {
	cobra.OnInitialize(onInit)

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
	rootCmd.Flags().StringP("loglevel", "l", "INFO", "")
	rootCmd.PersistentFlags().StringP("key", "k", "", "utc key file json content")
	rootCmd.PersistentFlags().StringP("secret", "s", "", "password to decrypt key file")

	// mine command initialize
	mineCmd.Flags().StringP("client-id", "c", "", "unique client id assigned to miner (required)")

	// stake command initialize
	stakeCmd.Flags().Int64("amount", 10, "amount to stake (default: wei)")

	// withdraw command initialize
	withdrawCmd.Flags().Int64("amount", 0, "amount to withdraw")

	// add commands and execute
	rootCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(stakeCmd)
	rootCmd.AddCommand(withdrawCmd)

	rootCmd.Execute()
}

func onInit() {
	loglevel := viper.GetString("loglevel")
	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		loglevel = logrus.InfoLevel.String()
		level, _ = logrus.ParseLevel(loglevel)
	}

	l := logrus.New()
	l.SetLevel(level)
	l.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano})

	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN != "" {
		sentryLevels := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		}
		sentryTags := map[string]string{
			"service": ServiceName,
			"version": Version,
		}
		sentryHook, err := logrussentry.NewAsyncWithTagsSentryHook(
			sentryDSN,
			sentryTags,
			sentryLevels,
		)
		sentryHook.StacktraceConfiguration.Enable = true
		sentryHook.Timeout = 5 * time.Second
		sentryHook.SetRelease(Version)

		if err != nil {
			l.Warning(err)
		} else {
			l.AddHook(sentryHook)
		}
	}

	cfg = service.LoadConfig()
	cfg.Logger = logrus.NewEntry(l)
	cfg.Name = ServiceName
	cfg.Version = Version
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

	go svc.Start()

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

	amount, err := cmd.Flags().GetInt64("amount")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cli.Register(context.Background(), amount)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("transcoder successfully registered")

}

func runWithdrawCommand(cmd *cobra.Command, args []string) {
	fmt.Println("run withdraw command")
	fmt.Printf("KEY=%s\n", viper.GetString("key"))
	fmt.Printf("SECRET=%s\n", viper.GetString("secret"))
}
