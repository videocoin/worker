package transcode

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/consul/api"
	bc "github.com/videocoin/cloud-pkg/bcops"
	"github.com/videocoin/cloud-pkg/stream"
	streamManager "github.com/videocoin/cloud-pkg/streamManager"

	"github.com/ethereum/go-ethereum/common"
)

// NewEth returns new eth object with new key
func NewEth(c *Config) *Eth {
	consul, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%s", os.Getenv("CONFIG_CONSUL_UI_SERVICE_HOST"), os.Getenv("CONFIG_CONSUL_UI_SERVICE_PORT_HTTP")),
	})

	if err != nil {
		panic(err)
	}

	keyPairs, _, err := consul.KV().List("config/dev/services/transcoder/keys", nil)
	if err != nil {
		panic(err)
	}

	key := keyPairs[rand.Intn(len(keyPairs)-1)]

	secret, _, err := consul.KV().Get(strings.Replace(key.Key, "keys", "secrets", 1), nil)
	if err != nil {
		panic(err)
	}

	rawKey, err := bc.LoadBcPrivKeys(string(key.Value), string(secret.Value), bc.FromMemory)
	if err != nil {
		panic(err)
	}

	return &Eth{
		rawKey: rawKey,
		kv:     key,
	}
}

func (e *Eth) connect(url, smca, streamAddress string) {
	var err error

	e.client, err = ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	managerAddress := common.HexToAddress(smca)

	e.sm, err = streamManager.NewManager(managerAddress, e.client)
	if err != nil {
		panic(err)
	}

	e.auth, err = bc.GetBCAuth(e.client, e.rawKey)
	if err != nil {
		panic(err)
	}

	e.si, err = stream.NewStream(common.HexToAddress(streamAddress), e.client)
	if err != nil {
		panic(err)
	}

	profiles, err := e.si.Getprofiles(nil)
	if err != nil {
		panic(err)
	}

	e.profiles = profiles
}
