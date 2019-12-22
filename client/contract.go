package client

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/videocoin/cloud-pkg/stakingManager"
	"github.com/videocoin/transcode/caller"
)

type StakingContract struct {
	instance *stakingManager.Staking
	common.Address
}

func NewStakingContract(managerAddress string, caller *caller.Caller) (*StakingContract, error) {
	addr := common.HexToAddress(managerAddress)
	instance, err := stakingManager.NewStaking(addr, caller.Client())
	if err != nil {
		return nil, err
	}

	return &StakingContract{instance, addr}, nil
}

func (m *StakingContract) IsTranscoder(addr common.Address) (bool, error) {
	isTranscoder, err := m.instance.IsTranscoder(&bind.CallOpts{}, addr)
	if err != nil {
		return false, err
	}

	return isTranscoder, nil
}
