package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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

func main() {
	cobra.OnInitialize(preInit)

	var rootCmd = &cobra.Command{
		Use: "",
	}

	var mineCmd = &cobra.Command{
		Use: "mine",
		Run: runMineCommand,
	}

	var stakeCmd = &cobra.Command{
		Use: "stake",
		Run: runStakeCommand,
	}

	var withdrawCmd = &cobra.Command{
		Use: "withdraw",
		Run: runWithdrawCommand,
	}

	viper.AutomaticEnv()

	// root command initialize
	rootCmd.Flags().StringP("loglevel", "l", "INFO", "")
	rootCmd.Flags().StringP("key", "k", "", "utc key file json content")
	rootCmd.Flags().StringP("secret", "s", "", "password to decrypt key file")

	rootCmd.MarkFlagRequired("loglevel")
	rootCmd.MarkFlagRequired("key")
	rootCmd.MarkFlagRequired("secret")

	viper.BindPFlags(rootCmd.Flags())

	// mine command initialize
	mineCmd.Flags().String("client-id", "", "unique client id assigned to miner (required)")
	mineCmd.MarkFlagRequired("client-id")
	viper.BindPFlag("client_id", mineCmd.Flags().Lookup("client-id"))

	// stake command initialize
	stakeCmd.Flags().Int64("amount", 10, "amount to stake (default: wei)")
	// stakeCmd.MarkFlagRequired("amount")

	// withdraw command initialize
	withdrawCmd.Flags().Int64("amount", 0, "amount to withdraw")
	withdrawCmd.MarkFlagRequired("amount")

	// add commands and execute
	rootCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(stakeCmd)
	rootCmd.AddCommand(withdrawCmd)

	rootCmd.Execute()
}

func preInit() {
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

	cfg.Key = viper.GetString("key")
	cfg.Secret = viper.GetString("secret")
	cfg.ClientID = viper.GetString("client_id")

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
	fmt.Println("run stake command")
	fmt.Printf("KEY=%s\n", viper.GetString("key"))
	fmt.Printf("SECRET=%s\n", viper.GetString("secret"))

	log := cfg.Logger

	cfg.Key = viper.GetString("key")
	cfg.Secret = viper.GetString("secret")

	cli, err := getTranscoderClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	amount, err := cmd.Flags().GetInt64("amount")
	if err != nil {
		log.Fatal(err.Error())
	}

	// register if it is not
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
