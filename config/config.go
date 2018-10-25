package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config default config for transcoder
type Config struct {
	LogLevel      string `required:"true" default:"debug" envconfig:"LOG_LEVEL" default:"DEBUG"`
	SQLURI        string `required:"true" envconfig:"SQL_URI" default:"root:password@/videocoin?parseTime=true"`
	MqURI         string `required:"true" envconfig:"MQ_URI" default:"amqp://guest:guest@127.0.0.1:5672"`
	ConsulAddr    string `required:"false" envconfig:"CONSUL_ADDR" default:"http://localhost:8500"`
	BaseStreamURL string `required:"true" envconfig:"BASE_STREAM_URL" default:"http://127.0.0.1:1935/hls/"`
}

var cfg Config
var once sync.Once

// Load initialize config
func Load() *Config {
	once.Do(func() {
		err := envconfig.Process("", &cfg)
		if err != nil {
			logrus.Fatalf("failed to load config: %s", err.Error())
		}
	})

	return &cfg
}
