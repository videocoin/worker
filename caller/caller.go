package caller

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

type Caller struct {
	client *ethclient.Client
	key    *keystore.Key
}

func NewCaller(jsonkey, pwd string, client *ethclient.Client) (*Caller, error) {
	key, err := keystore.DecryptKey([]byte(jsonkey), pwd)
	if err != nil {
		return nil, err
	}

	return &Caller{
		client: client,
		key:    key,
	}, nil
}

func (c *Caller) Addr() common.Address {
	return c.key.Address
}

func (c *Caller) Balance() (*big.Int, error) {
	return c.client.BalanceAt(context.Background(), c.key.Address, nil)
}

func (c *Caller) Client() *ethclient.Client {
	return c.client
}

func (c *Caller) TransactOpts(gasLimit int) *bind.TransactOpts {
	opts := bind.NewKeyedTransactor(c.key.PrivateKey)
	if gasLimit == 0 {
		opts.GasLimit = uint64(3000000)
	}

	return opts
}

func (c *Caller) WaitMinedAndCheck(tx *types.Transaction) error {
	cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(cancelCtx, c.client, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction %s failed", tx.Hash().String())
	}

	return nil
}
