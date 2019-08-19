package transcode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

	rawKey, err := bc.LoadBcPrivKeys(`{"address":"263287e4debff1cd61d75127794d837ba4932928","crypto":{"cipher":"aes-128-ctr","ciphertext":"a805d15481f97b6fa2500b7a6876dddbd4a527f0b240c2890a67cb31f0edc3b6","cipherparams":{"iv":"6dfc7c3b3a280a175f131323ec6d2670"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"7f3f92d81110ea68b896b6ba5db2164d514822dc0263b8bdb2747a3572221ca1"},"mac":"7ef9d43e4068608617d3db2b590ae3e89a0a5766616d81f7fa8c94205a98c77d"},"id":"0ca0e295-b4fd-451f-8ba9-0d5f8bd38108","version":3}`, "efbzjzcqmg", bc.FromMemory)
	if err != nil {
		panic(err)
	}

	fmt.Println(rawKey)

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
