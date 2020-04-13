package cryptoinfo

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"

	staking "github.com/videocoin/go-staking"
	"github.com/videocoin/transcode/caller"
)

type CryptoInfo struct {
	caller *caller.Caller
	client *staking.Client
}

func NewCryptoInfo(caller *caller.Caller, addr string) (*CryptoInfo, error) {
	client, err := staking.NewClient(caller.NatClient(), common.HexToAddress(addr))
	if err != nil {
		return nil, err
	}

	return &CryptoInfo{
		caller: caller,
		client: client,
	}, nil
}

func (ci *CryptoInfo) GetInfo() (map[string]interface{}, []byte, error) {
	info := map[string]interface{}{}

	stake, err := ci.client.GetTranscoderStake(context.Background(), ci.caller.Addr())
	if err == nil {
		info["stake"] = stake.String()
	}

	balance, err := ci.caller.Balance()
	if err == nil {
		info["balance"] = balance.String()
	}

	info["address"] = ci.caller.Addr().String()

	b, err := json.Marshal(info)
	if err != nil {
		return nil, nil, err
	}

	return info, b, nil
}
