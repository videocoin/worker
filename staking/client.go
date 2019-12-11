package staking

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/videocoin/c/caller"
	"github.com/videocoin/cloud-pkg/stakingManager"
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

func (c *TranscoderClient) Register(ctx context.Context, opts *bind.TransactOpts) error {
	addr := opts.From

	isTranscoder, err := c.IsTranscoder(addr)
	if err != nil {
		return err
	}

	if isTranscoder {
		return fmt.Errorf("address: %s is already a transcoder", addr)
	}

	tx, err := c.instance.TranscoderRegister(opts)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.client, tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *TranscoderClient) SelfStake(ctx context.Context, opts *bind.TransactOpts) error {
	addrTran := opts.From

	isTranscoder, err := c.IsTranscoder(addrTran)
	if err != nil {
		return err
	}

	if !isTranscoder {
		return fmt.Errorf("address: %s is not a transcoder", addrTran)
	}

	tx, err := c.instance.SelfStake(opts)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.client, tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *TranscoderClient) WithdrawalProposal(ctx context.Context, opts *bind.TransactOpts) error {
	addrTran := opts.From

	isTranscoder, err := c.IsTranscoder(addrTran)
	if err != nil {
		return err
	}

	if !isTranscoder {
		return fmt.Errorf("address: %s is not a transcoder", addrTran)
	}

	tx, err := c.instance.WithdrawalTranscoderProposal(opts)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.client, tx)
	if err != nil {
		return err
	}

	return nil
}

func (c *TranscoderClient) WithdrawStake(ctx context.Context, amount *big.Int, opts *bind.TransactOpts) error {
	addrTran := opts.From

	isTranscoder, err := c.IsTranscoder(addrTran)
	if err != nil {
		return err
	}

	if !isTranscoder {
		return fmt.Errorf("address: %s is not a transcoder", addrTran)
	}

	tx, err := c.instance.WithdrawTranscoderStake(opts, amount)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, c.client, tx)
	if err != nil {
		return err
	}

	return nil
}
