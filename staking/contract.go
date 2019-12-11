package staking

import (
	"context"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	staking "github.com/videocoin/cloud-pkg/stakingManager"
)

type StakingContract struct {
	instance *staking.Staking
	common.Address
}

func NewStakingContract(managerAddress string, caller *caller.Caller) (*StakingContract, error) {
	addr := common.HexToAddress(managerAddress)
	instance, err := staking.NewStaking(addr, client)
	if err != nil {
		return nil, err
	}

	return &StakingContract{instance, addr}, nil
}

// returns true if address is registerred as transcoder
func (m *StakingContract) IsTranscoder(addr common.Address) (bool, error) {
	isTranscoder, err := m.instance.IsTranscoder(&bind.CallOpts{}, addr)
	if err != nil {
		return false, err
	}

	return isTranscoder, nil
}