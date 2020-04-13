package caller

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Caller struct {
	ethClient *ethclient.Client
	natClient *ethclient.Client
	key       *keystore.Key
}

func NewCaller(jsonkey, pwd string, natClient *ethclient.Client, ethClient *ethclient.Client) (*Caller, error) {
	key, err := keystore.DecryptKey([]byte(jsonkey), pwd)
	if err != nil {
		return nil, err
	}

	return &Caller{
		natClient: natClient,
		key:       key,
	}, nil
}

func (c *Caller) PrivateKey() *ecdsa.PrivateKey {
	return c.key.PrivateKey
}

func (c *Caller) Addr() common.Address {
	return c.key.Address
}

func (c *Caller) Balance() (*big.Int, error) {
	return c.natClient.BalanceAt(context.Background(), c.key.Address, nil)
}

func (c *Caller) NatClient() *ethclient.Client {
	return c.natClient
}

//func (c *Caller) EthClient() *ethclient.Client {
//	return c.ethClient
//}

func (c *Caller) Opts(amount *big.Int) *bind.TransactOpts {
	gasPrice, _ := c.natClient.SuggestGasPrice(context.Background())

	opts := bind.NewKeyedTransactor(c.key.PrivateKey)
	opts.Nonce = nil
	opts.Value = amount
	opts.GasPrice = gasPrice
	opts.GasLimit = uint64(8000000)

	return opts
}

func (c *Caller) WaitMinedAndCheck(tx *types.Transaction) error {
	cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(cancelCtx, c.natClient, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction %s failed", tx.Hash().String())
	}

	return nil
}
