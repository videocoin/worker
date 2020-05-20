package stream

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/videocoin/go-protocol/streams"
	"github.com/videocoin/worker/caller"
)

type Contract struct {
	instance *streams.Stream
	common.Address
}

func NewContract(streamAddr string, caller *caller.Caller) (*Contract, error) {
	addr := common.HexToAddress(streamAddr)
	instance, err := streams.NewStream(addr, caller.NatClient())
	if err != nil {
		return nil, err
	}

	return &Contract{instance, addr}, nil
}
