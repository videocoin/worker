package transcode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
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
		Address: c.ConsulAddress,
		HttpAuth: &api.HttpBasicAuth{
			Username: c.ConsulUsername,
			Password: c.ConsulPassword,
		},
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

	key.Value = decrypt(key.Value, c.DecryptionKey)
	secret.Value = decrypt(secret.Value, c.DecryptionKey)

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

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return plaintext
}
