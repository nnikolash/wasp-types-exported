package chainutil

import (
	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
)

func GetAccountBalance(ch chain.ChainCore, agentID isc.AgentID) (*isc.Assets, error) {
	params := codec.MakeDict(map[string]interface{}{
		accounts.ParamAgentID: codec.EncodeAgentID(agentID),
	})
	ret, err := CallView(mustLatestState(ch), ch, accounts.Contract.Hname(), accounts.ViewBalance.Hname(), params)
	if err != nil {
		return nil, err
	}
	return isc.AssetsFromDict(ret)
}

func mustLatestState(ch chain.ChainCore) state.State {
	latestState, err := ch.LatestState(chain.ActiveOrCommittedState)
	if err != nil {
		panic(err)
	}
	return latestState
}
