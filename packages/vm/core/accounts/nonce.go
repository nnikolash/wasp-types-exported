package accounts

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
)

func NonceKey(callerAgentID isc.AgentID, chainID isc.ChainID) kv.Key {
	return KeyNonce + AccountKey(callerAgentID, chainID)
}

// Nonce returns the "total request count" for an account (it's the AccountNonce that is expected in the next request)
func AccountNonce(state kv.KVStoreReader, callerAgentID isc.AgentID, chainID isc.ChainID) uint64 {
	if callerAgentID.Kind() == isc.AgentIDKindEthereumAddress {
		panic("to get EVM nonce, call EVM contract")
	}
	data := state.Get(NonceKey(callerAgentID, chainID))
	if data == nil {
		return 0
	}
	return codec.MustDecodeUint64(data) + 1
}

func IncrementNonce(state kv.KVStore, callerAgentID isc.AgentID, chainID isc.ChainID) {
	if callerAgentID.Kind() == isc.AgentIDKindEthereumAddress {
		// don't update EVM nonces
		return
	}
	next := AccountNonce(state, callerAgentID, chainID)
	state.Set(NonceKey(callerAgentID, chainID), codec.EncodeUint64(next))
}
