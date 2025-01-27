package accounts

import (
	"math/big"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/util"
)

type (
	getBaseTokensFn             func(state kv.KVStoreReader, accountKey kv.Key) uint64
	GetBaseTokensFullDecimalsFn func(state kv.KVStoreReader, accountKey kv.Key) *big.Int
	setBaseTokensFullDecimalsFn func(state kv.KVStore, accountKey kv.Key, amount *big.Int)
)

func getBaseTokens(v isc.SchemaVersion) getBaseTokensFn {
	switch v {
	case 0:
		return getBaseTokensDEPRECATED
	default:
		return getBaseTokensNEW
	}
}

func GetBaseTokensFullDecimals(v isc.SchemaVersion) GetBaseTokensFullDecimalsFn {
	switch v {
	case 0:
		return getBaseTokensFullDecimalsDEPRECATED
	default:
		return getBaseTokensFullDecimalsNEW
	}
}

func setBaseTokensFullDecimals(v isc.SchemaVersion) setBaseTokensFullDecimalsFn {
	switch v {
	case 0:
		return setBaseTokensFullDecimalsDEPRECATED
	default:
		return setBaseTokensFullDecimalsNEW
	}
}

// -------------------------------------------------------------------------------

func BaseTokensKey(accountKey kv.Key) kv.Key {
	return PrefixBaseTokens + accountKey
}

func getBaseTokensFullDecimalsNEW(state kv.KVStoreReader, accountKey kv.Key) *big.Int {
	return codec.MustDecodeBigIntAbs(state.Get(BaseTokensKey(accountKey)), big.NewInt(0))
}

func setBaseTokensFullDecimalsNEW(state kv.KVStore, accountKey kv.Key, amount *big.Int) {
	state.Set(BaseTokensKey(accountKey), codec.EncodeBigIntAbs(amount))
}

func getBaseTokensNEW(state kv.KVStoreReader, accountKey kv.Key) uint64 {
	amount := getBaseTokensFullDecimalsNEW(state, accountKey)
	// convert from 18 decimals, discard the remainder
	convertedAmount, _ := util.EthereumDecimalsToBaseTokenDecimals(amount, parameters.L1().BaseToken.Decimals)
	return convertedAmount
}

func AdjustAccountBaseTokens(v isc.SchemaVersion, state kv.KVStore, account isc.AgentID, adjustment int64, chainID isc.ChainID) {
	switch {
	case adjustment > 0:
		CreditToAccount(v, state, account, isc.NewAssets(uint64(adjustment), nil), chainID)
	case adjustment < 0:
		DebitFromAccount(v, state, account, isc.NewAssets(uint64(-adjustment), nil), chainID)
	}
}

func GetBaseTokensBalance(v isc.SchemaVersion, state kv.KVStoreReader, agentID isc.AgentID, chainID isc.ChainID) uint64 {
	return getBaseTokens(v)(state, AccountKey(agentID, chainID))
}

func GetBaseTokensBalanceFullDecimals(v isc.SchemaVersion, state kv.KVStoreReader, agentID isc.AgentID, chainID isc.ChainID) *big.Int {
	return GetBaseTokensFullDecimals(v)(state, AccountKey(agentID, chainID))
}
