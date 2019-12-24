package stream

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/videocoin/transcode/caller"
)

type StreamClient struct {
	StreamContract
	caller.Caller
}

func NewStreamClient(addr string, caller *caller.Caller) (*StreamClient, error) {
	contract, err := NewStreamContract(addr, caller)
	if err != nil {
		return nil, err
	}

	return &StreamClient{*contract, *caller}, nil
}

func (sc *StreamClient) GetProfiles() ([]*big.Int, error) {
	profiles, err := sc.instance.Getprofiles(nil)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (sc *StreamClient) GetInChunks() ([]*big.Int, error) {
	chunks, err := sc.instance.GetInChunks(nil)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func (sc *StreamClient) SubmitProof(chunkID, outChunkID *big.Int, profileID *big.Int) (*types.Transaction, error) {
	tx, err := sc.instance.SubmitProof(sc.TransactOpts(big.NewInt(0), 0), profileID, chunkID, big.NewInt(0), outChunkID)
	if err != nil {
		return nil, err
	}

	_, err = bind.WaitMined(context.Background(), sc.Client(), tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
