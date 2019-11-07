package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	logrussentry "github.com/evalphobia/logrus_sentry"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/transcode/service"
)

var (
	ServiceName string = "transcoder"
	Version     string = "dev"
)

func main() {
	loglevel := os.Getenv("LOGLEVEL")
	if loglevel == "" {
		loglevel = logrus.InfoLevel.String()
	}

	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		loglevel = logrus.InfoLevel.String()
		level, _ = logrus.ParseLevel(loglevel)
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})

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
			logrus.Warning(err)
		} else {
			logrus.AddHook(sentryHook)
		}
	}

	log := logrus.NewEntry(logrus.New())

	closer, err := tracer.NewTracer(ServiceName)
	if err != nil {
		log.Info(err.Error())
	} else {
		defer closer.Close()
	}

	cfg := &service.Config{
		Name:    ServiceName,
		Version: Version,
	}

	err = envconfig.Process(ServiceName, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg.Logger = log

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
