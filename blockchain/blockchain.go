package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	bc "github.com/videocoin/cloud-pkg/bcops"
	"github.com/videocoin/cloud-pkg/stream"
	// streamManager "github.com/videocoin/cloud-pkg/streamManager"
)

type Config struct {
	URL    string
	SMCA   string
	Key    string
	Secret string
}

type Client struct {
	cli *ethclient.Client
	// sm   *streamManager.Manager
	auth *bind.TransactOpts
}

func Dial(cfg *Config) (*Client, error) {
	var err error
	eth := &Client{}

	rawKey, err := bc.LoadBcPrivKeys(cfg.Key, cfg.Secret, bc.FromMemory)
	if err != nil {
		return nil, err
	}

	eth.cli, err = ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	// managerAddress := ethcommon.HexToAddress(cfg.SMCA)
	// eth.sm, err = streamManager.NewManager(managerAddress, eth.cli)
	// if err != nil {
	// 	return nil, err
	// }

	eth.auth, err = bc.GetBCAuth(eth.cli, rawKey)
	if err != nil {
		return nil, err
	}

	return eth, nil
}

func (c *Client) EthClient() *ethclient.Client {
	return c.cli
}

func (c *Client) EthAuth() *bind.TransactOpts {
	return c.auth
}

type StreamContract struct {
	cli    *ethclient.Client
	stream *stream.Stream
	auth   *bind.TransactOpts
}

func NewStreamContract(addr string, cli *ethclient.Client, auth *bind.TransactOpts) (*StreamContract, error) {
	var err error

	stream, err := stream.NewStream(ethcommon.HexToAddress(addr), cli)
	if err != nil {
		return nil, err
	}

	return &StreamContract{
		cli:    cli,
		stream: stream,
		auth:   auth,
	}, nil
}

func (sc *StreamContract) GetProfiles() ([]*big.Int, error) {
	profiles, err := sc.stream.Getprofiles(nil)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (sc *StreamContract) GetInChunks() ([]*big.Int, error) {
	chunks, err := sc.stream.GetInChunks(nil)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func (sc *StreamContract) SubmitProof(chunkID, outChunkID *big.Int, profileID *big.Int) (*ethtypes.Transaction, error) {
	tx, err := sc.stream.SubmitProof(sc.auth, profileID, chunkID, big.NewInt(0), outChunkID)
	if err != nil {
		return nil, err
	}

	err = sc.waitMinedAndCheck(tx)
	if err != nil {
		return tx, fmt.Errorf("failed to wait mined after submit proof: %s", err.Error())
	}

	return tx, nil
}

func (sc *StreamContract) waitMinedAndCheck(tx *ethtypes.Transaction) error {
	cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(cancelCtx, sc.cli, tx)
	if err != nil {
		return err
	}

	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction %s failed", tx.Hash().String())
	}

	return nil
}
