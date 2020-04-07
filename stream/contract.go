package stream

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/videocoin/cloud-pkg/stream"
	"github.com/videocoin/transcode/caller"
)

type Contract struct {
	instance *stream.Stream
	common.Address
}

func NewContract(streamAddr string, caller *caller.Caller) (*Contract, error) {
	addr := common.HexToAddress(streamAddr)
	instance, err := stream.NewStream(addr, caller.EthClient())
	if err != nil {
		return nil, err
	}

	return &Contract{instance, addr}, nil
}
