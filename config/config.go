package config

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/datastore"
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
	OutputDir     string `required:"true" envconfig:"OUTPUT_DIR" default:"/opt/mnt/" description:"Mount point for GCSFUSE"`
}

var cfg Config
var once sync.Once

// Load initialize config
func Load(loc string) *Config {
	switch loc {
	case "local":
		once.Do(func() {
			err := envconfig.Process("", &cfg)
			if err != nil {
				logrus.Fatalf("failed to load config: %s", err.Error())
			}
		})
		break
	// requires PROJECT_ID environment variable
	case "remote":
		ctx := context.Background()
		client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
		if err != nil {
			logrus.Fatalf("failed to create new client: %s", err)
		}

		key := datastore.NameKey("config", "transcoder", nil)
		err = client.Get(ctx, key, &cfg)
		if err != nil {
			logrus.Fatalf("failed to get namekey: %s", err)
		}

		break

	default:
	}

	return &cfg
}
