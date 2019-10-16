package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	env := os.Getenv("CLUSTER_ENV")
	outPath := os.Getenv("OUTPUT_PATH")
	if outPath == "" {
		outPath = "/env/init.env"
	}
	consul, err := consulapi.NewClient(&consulapi.Config{
		Address: fmt.Sprintf(
			"%s:%s",
			os.Getenv("CONFIG_CONSUL_UI_SERVICE_HOST"),
			os.Getenv("CONFIG_CONSUL_UI_SERVICE_PORT_HTTP"),
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	keyPath := fmt.Sprintf("config/%s/services/transcoder/keys", env)
	keyPairs, _, err := consul.KV().List(keyPath, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(keyPairs) > 0 {
		rand.Seed(time.Now().UnixNano())
		key := keyPairs[rand.Intn(len(keyPairs)-1)]

		k, _, err := consul.KV().Get(key.Key, nil)
		if err != nil {
			log.Fatal(err)
		}

		parts := strings.Split(key.Key, "/")
		if len(parts) > 0 {
			secretBase64 := parts[len(parts)-1]
			secret, err := base64.StdEncoding.DecodeString(secretBase64)
			if err != nil {
				log.Fatal(err)
			}

			keyValue := strings.Replace(string(k.Value), "'", "\"", -1)
			content := fmt.Sprintf("export KEY='%s'\nexport SECRET=%s\n", keyValue, secret)
			err = ioutil.WriteFile(outPath, []byte(content), 0777)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(content)
		}
	}
}
