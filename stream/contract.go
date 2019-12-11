package stream

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/videocoin/cloud-pkg/stream"
	"github.com/videocoin/transcode/caller"
)

type StreamContract struct {
	instance *stream.Stream
	common.Address
}

func NewStreamContract(streamAddr string, caller *caller.Caller) (*StreamContract, error) {
	addr := common.HexToAddress(streamAddr)
	instance, err := stream.NewStream(addr, caller.Client())
	if err != nil {
		return nil, err
	}

	return &StreamContract{instance, addr}, nil
}
