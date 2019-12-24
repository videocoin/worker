package client

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/videocoin/transcode/caller"
)

var (
	ErrAlreadyRegistered = errors.New("transcoder address is already registered")
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

func (c *TranscoderClient) isRegistered(addr common.Address) error {
	isTranscoder, err := c.IsTranscoder(addr)
	if err != nil {
		return err
	}

	if isTranscoder {
		return ErrAlreadyRegistered
	}

	return nil
}

func (c *TranscoderClient) Register(ctx context.Context, amount *big.Int) error {
	opts := c.TransactOpts(amount, 0)

	err := c.isRegistered(opts.From)
	if err != nil {
		return err
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

func (c *TranscoderClient) SelfStake(ctx context.Context, amount *big.Int) error {
	opts := c.TransactOpts(amount, 0)

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
	opts := c.TransactOpts(big.NewInt(0), 0)

	err := c.isRegistered(opts.From)
	if err != nil && err != ErrAlreadyRegistered {
		return err
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
	opts := c.TransactOpts(big.NewInt(0), 0)

	err := c.isRegistered(opts.From)
	if err != nil && err != ErrAlreadyRegistered {
		return err
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

func (c *TranscoderClient) GetStake() (*big.Int, error) {
	opts := c.TransactOpts(big.NewInt(0), 0)

	err := c.isRegistered(opts.From)
	if err != nil && err != ErrAlreadyRegistered {
		return nil, err
	}

	return c.instance.GetTranscodersStake(&bind.CallOpts{}, opts.From)
}
