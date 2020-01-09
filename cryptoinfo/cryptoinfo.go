package cryptoinfo

import (
	"encoding/json"

	"github.com/videocoin/transcode/caller"
	"github.com/videocoin/transcode/client"
)

type CryptoInfo struct {
	caller  *caller.Caller
	tclient *client.TranscoderClient
}

func NewCryptoInfo(addr string, caller *caller.Caller) (*CryptoInfo, error) {
	tcli, err := client.NewTranscoderClient(addr, caller)
	if err != nil {
		return nil, err
	}

	return &CryptoInfo{
		caller:  caller,
		tclient: tcli,
	}, nil
}

func (ci *CryptoInfo) GetInfo() (map[string]interface{}, []byte, error) {
	info := map[string]interface{}{}

	stake, err := ci.tclient.GetStake()
	if err == nil {
		info["self_stake"] = stake.String()
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
