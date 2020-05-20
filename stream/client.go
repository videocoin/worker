package stream

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/videocoin/worker/caller"
)

type Client struct {
	Contract
	caller.Caller
}

func NewClient(addr string, caller *caller.Caller) (*Client, error) {
	contract, err := NewContract(addr, caller)
	if err != nil {
		return nil, err
	}

	return &Client{*contract, *caller}, nil
}

func (sc *Client) GetProfiles() ([]*big.Int, error) {
	profiles, err := sc.instance.Getprofiles(nil)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (sc *Client) GetInChunks() ([]*big.Int, error) {
	chunks, err := sc.instance.GetInChunks(nil)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func (sc *Client) SubmitProof(chunkID, outChunkID *big.Int, profileID *big.Int) (*types.Transaction, error) {
	tx, err := sc.instance.SubmitProof(sc.Opts(nil), profileID, chunkID, big.NewInt(0), outChunkID)
	if err != nil {
		return nil, err
	}

	_, err = bind.WaitMined(context.Background(), sc.NatClient(), tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
