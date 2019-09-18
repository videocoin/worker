package transcode

import (
	"math/big"
	"os/exec"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/grafov/m3u8"
	"github.com/hashicorp/consul/api"
	"github.com/videocoin/cloud-pkg/stream"
	streamManager "github.com/videocoin/cloud-pkg/streamManager"
)

// Byte constants
const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

type (
	// Task used to queue up chunk work
	Task struct {
		ID              string
		Bitrate         uint32
		InputChunkName  string
		OutputChunkName string
		ChunksDir       string
		StreamAddress   string
		StreamID        *big.Int
		InputID         *big.Int
		OutputID        *big.Int
		cmd             *exec.Cmd
		Wallet          common.Address
		Playlist        *m3u8.MediaPlaylist
	}

	// Eth Used for all eth interactions
	Eth struct {
		sm       *streamManager.Manager
		si       *stream.Stream
		kv       *api.KVPair
		rawKey   *keystore.Key
		auth     *bind.TransactOpts
		client   *ethclient.Client
		profiles []*big.Int
	}
)
