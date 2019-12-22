package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/videocoin/transcode/caller"
)

type TranscoderClient struct {
	StakingContract
	caller.Caller
}

func NewTranscoderClient(managerAddr string, caller *caller.Caller) (*TranscoderClient, error) {
	contract, err := NewStakingContract(managerAddr, caller)
	if err != nil {
		return nil, err
	}

	return &TranscoderClient{*contract, *caller}, nil
}

func (c *TranscoderClient) Register(ctx context.Context) error {
	opts := c.TransactOpts(0)

	isTranscoder, err := c.IsTranscoder(opts.From)
	if err != nil {
		return err
	}

	if isTranscoder {
		return fmt.Errorf("address: %s is already a transcoder", opts.From)
	}

	tx, err := c.instance.TranscoderRegister(opts)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.Client(), tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *TranscoderClient) SelfStake(ctx context.Context) error {
	opts := c.TransactOpts(0)

	isTranscoder, err := c.IsTranscoder(opts.From)
	if err != nil {
		return err
	}

	if !isTranscoder {
		return fmt.Errorf("address: %s is not a transcoder", opts.From)
	}

	tx, err := c.instance.SelfStake(opts)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.Client(), tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *TranscoderClient) WithdrawalProposal(ctx context.Context) error {
	opts := c.TransactOpts(0)

	isTranscoder, err := c.IsTranscoder(opts.From)
	if err != nil {
		return err
	}

	if !isTranscoder {
		return fmt.Errorf("address: %s is not a transcoder", opts.From)
	}

	tx, err := c.instance.WithdrawalTranscoderProposal(opts)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.Client(), tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *TranscoderClient) WithdrawStake(ctx context.Context, amount *big.Int) error {
	opts := c.TransactOpts(0)

	isTranscoder, err := c.IsTranscoder(opts.From)
	if err != nil {
		return err
	}

	if !isTranscoder {
		return fmt.Errorf("address: %s is not a transcoder", opts.From)
	}

	tx, err := c.instance.WithdrawTranscoderStake(opts, amount)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.Client(), tx)
	if err != nil {
		return err
	}

	return nil
}
