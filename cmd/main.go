package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	logrussentry "github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/transcode/service"
)

var (
	ServiceName string = "transcoder"
	Version     string = "dev"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "transcoder",
	}

	var mineCmd = &cobra.Command{
		Use: "mine",
		Run: runMineCommand,
	}

	var registerCmd = &cobra.Command{
		Use: "register",
		Run: runRegisterCommand,
	}

	var stakeCmd = &cobra.Command{
		Use: "stake",
		Run: runStakeCommand,
	}

	var withDrawCmd = &cobra.Command{
		Use: "withdraw",
		Run: runWithdrawCommand,
	}

	viper.AutomaticEnv()

	// root command initialize
	rootCmd.Flags().StringP("loglevel", "l", "INFO", "")
	rootCmd.Flags().StringP("key", "k", "", "")
	rootCmd.Flags().StringP("secret", "s", "", "")

	rootCmd.MarkFlagRequired("loglevel")
	rootCmd.MarkFlagRequired("key")
	rootCmd.MarkFlagRequired("secret")

	viper.BindPFlag("loglevel", rootCmd.Flags().Lookup("loglevel"))
	viper.BindPFlag("key", rootCmd.Flags().Lookup("key"))
	viper.BindPFlag("secret", rootCmd.Flags().Lookup("secret"))

	// mine command initialize
	mineCmd.Flags().String("output-dir", "/tmp", "")
	mineCmd.Flags().String("client-id", "", "")
	mineCmd.Flags().String("dispatcher-addr", "d.dev.videocoin.network:5008", "")
	mineCmd.Flags().String("blockchain-url", "https://dev1:D6msEL93LJT5RaPk@rpc.dev.kili.videocoin.network", "")
	mineCmd.Flags().String("syncer-url", "https://dev.videocoin.network/api/v1/sync", "")

	viper.BindPFlag("output_dir", mineCmd.Flags().Lookup("output-dir"))
	viper.BindPFlag("client_id", mineCmd.Flags().Lookup("client-id"))

	viper.BindPFlag("dispatcher_addr", mineCmd.Flags().Lookup("dispatcher-addr"))
	viper.BindPFlag("blockchain_url", mineCmd.Flags().Lookup("blockchain-url"))
	viper.BindPFlag("syncer_url", mineCmd.Flags().Lookup("syncer-url"))

	// register command initialize
	// stake command initialize
	// withdraw command initialize

	// add commands and execute
	rootCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(stakeCmd)
	rootCmd.AddCommand(withDrawCmd)

	rootCmd.Execute()
}

func runMineCommand(cmd *cobra.Command, args []string) {
	loglevel := viper.GetString("loglevel")
	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		loglevel = logrus.InfoLevel.String()
		level, _ = logrus.ParseLevel(loglevel)
	}

	l := logrus.New()
	l.SetLevel(level)
	l.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano})
	log := logrus.NewEntry(l)

	// logrus.SetLevel(level)
	// logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})

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

	closer, err := tracer.NewTracer(ServiceName)
	if err != nil {
		log.Info(err.Error())
	} else {
		defer closer.Close()
	}

	cfg := &service.Config{
		Name:    ServiceName,
		Version: Version,
		Logger:  log,
	}

	cfg.Key = viper.GetString("key")
	cfg.Secret = viper.GetString("secret")
	cfg.ClientID = viper.GetString("client_id")

	cfg.OutputDir = viper.GetString("output_dir")
	cfg.DispatcherRPCAddr = viper.GetString("dispatcher_addr")
	cfg.RPCNodeURL = viper.GetString("blockchain_url")
	cfg.SyncerURL = viper.GetString("syncer_url")

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

func runRegisterCommand(cmd *cobra.Command, args []string) {
	fmt.Println("run register command")
	fmt.Printf("KEY=%s\n", viper.GetString("key"))
	fmt.Printf("SECRET=%s\n", viper.GetString("secret"))
}

func runStakeCommand(cmd *cobra.Command, args []string) {
	fmt.Println("run stake command")
	fmt.Printf("KEY=%s\n", viper.GetString("key"))
	fmt.Printf("SECRET=%s\n", viper.GetString("secret"))
}

func runWithdrawCommand(cmd *cobra.Command, args []string) {
	fmt.Println("run withdraw command")
	fmt.Printf("KEY=%s\n", viper.GetString("key"))
	fmt.Printf("SECRET=%s\n", viper.GetString("secret"))
}
