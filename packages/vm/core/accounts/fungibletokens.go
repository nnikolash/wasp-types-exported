package accounts

import (
	"fmt"
	"math/big"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/util"
)

// CreditToAccount brings new funds to the on chain ledger
// NOTE: this function does not take NFTs into account
func CreditToAccount(v isc.SchemaVersion, state kv.KVStore, agentID isc.AgentID, assets *isc.Assets, chainID isc.ChainID) {
	if assets == nil || assets.IsEmpty() {
		return
	}
	creditToAccount(v, state, AccountKey(agentID, chainID), assets)
	creditToAccount(v, state, L2TotalsAccount, assets)
	touchAccount(state, agentID, chainID)
}

// creditToAccount adds assets to the internal account map
// NOTE: this function does not take NFTs into account
func creditToAccount(v isc.SchemaVersion, state kv.KVStore, accountKey kv.Key, assets *isc.Assets) {
	if assets == nil || assets.IsEmpty() {
		return
	}

	if assets.BaseTokens > 0 {
		incomingTokensFullDecimals := util.BaseTokensDecimalsToEthereumDecimals(assets.BaseTokens, parameters.L1().BaseToken.Decimals)
		creditToAccountFullDecimals(v, state, accountKey, incomingTokensFullDecimals)
	}
	for _, nt := range assets.NativeTokens {
		if nt.Amount.Sign() == 0 {
			continue
		}
		if nt.Amount.Sign() < 0 {
			panic(ErrBadAmount)
		}
		balance := getNativeTokenAmount(state, accountKey, nt.ID)
		balance.Add(balance, nt.Amount)
		if balance.Cmp(util.MaxUint256) > 0 {
			panic(ErrOverflow)
		}
		setNativeTokenAmount(state, accountKey, nt.ID, balance)
	}
}

func CreditToAccountFullDecimals(v isc.SchemaVersion, state kv.KVStore, agentID isc.AgentID, amount *big.Int, chainID isc.ChainID) {
	if !util.IsPositiveBigInt(amount) {
		return
	}
	creditToAccountFullDecimals(v, state, AccountKey(agentID, chainID), amount)
	creditToAccountFullDecimals(v, state, L2TotalsAccount, amount)
	touchAccount(state, agentID, chainID)
}

// creditToAccountFullDecimals adds assets to the internal account map
func creditToAccountFullDecimals(v isc.SchemaVersion, state kv.KVStore, accountKey kv.Key, amount *big.Int) {
	setBaseTokensFullDecimals(v)(state, accountKey, new(big.Int).Add(GetBaseTokensFullDecimals(v)(state, accountKey), amount))
}

// DebitFromAccount takes out assets balance the on chain ledger. If not enough it panics
// NOTE: this function does not take NFTs into account
func DebitFromAccount(v isc.SchemaVersion, state kv.KVStore, agentID isc.AgentID, assets *isc.Assets, chainID isc.ChainID) {
	if assets == nil || assets.IsEmpty() {
		return
	}
	if !debitFromAccount(v, state, AccountKey(agentID, chainID), assets) {
		panic(fmt.Errorf("cannot debit (%s) from %s: %w", assets, agentID, ErrNotEnoughFunds))
	}
	if !debitFromAccount(v, state, L2TotalsAccount, assets) {
		panic("debitFromAccount: inconsistent ledger state")
	}
	touchAccount(state, agentID, chainID)
}

// debitFromAccount debits assets from the internal accounts map
// NOTE: this function does not take NFTs into account
func debitFromAccount(v isc.SchemaVersion, state kv.KVStore, accountKey kv.Key, assets *isc.Assets) bool {
	if assets == nil || assets.IsEmpty() {
		return true
	}

	// first check, then mutate
	mutateBaseTokens := false

	baseTokensToDebit := util.BaseTokensDecimalsToEthereumDecimals(assets.BaseTokens, parameters.L1().BaseToken.Decimals)
	var baseTokensToSet *big.Int
	if assets.BaseTokens > 0 {
		balance := GetBaseTokensFullDecimals(v)(state, accountKey)
		if baseTokensToDebit.Cmp(balance) > 0 {
			return false
		}
		mutateBaseTokens = true
		baseTokensToSet = new(big.Int).Sub(balance, baseTokensToDebit)
	}

	nativeTokensMutations := isc.NewEmptyAssets()
	for _, nt := range assets.NativeTokens {
		if nt.Amount.Sign() == 0 {
			continue
		}
		if nt.Amount.Sign() < 0 {
			panic(ErrBadAmount)
		}
		balance := getNativeTokenAmount(state, accountKey, nt.ID)
		balance = balance.Sub(balance, nt.Amount)
		if balance.Sign() < 0 {
			return false
		}
		nativeTokensMutations.AddNativeTokens(nt.ID, balance)
	}

	if mutateBaseTokens {
		setBaseTokensFullDecimals(v)(state, accountKey, baseTokensToSet)
	}
	for _, nt := range nativeTokensMutations.NativeTokens {
		setNativeTokenAmount(state, accountKey, nt.ID, nt.Amount)
	}
	return true
}

// DebitFromAccountFullDecimals removes the amount from the chain ledger. If not enough it panics
func DebitFromAccountFullDecimals(v isc.SchemaVersion, state kv.KVStore, agentID isc.AgentID, amount *big.Int, chainID isc.ChainID) {
	if !util.IsPositiveBigInt(amount) {
		return
	}
	if !debitFromAccountFullDecimals(v, state, AccountKey(agentID, chainID), amount) {
		panic(fmt.Errorf("cannot debit (%s) from %s: %w", amount.String(), agentID, ErrNotEnoughFunds))
	}

	if !debitFromAccountFullDecimals(v, state, L2TotalsAccount, amount) {
		panic("debitFromAccount: inconsistent ledger state")
	}
	touchAccount(state, agentID, chainID)
}

// debitFromAccountFullDecimals debits the amount from the internal accounts map
func debitFromAccountFullDecimals(v isc.SchemaVersion, state kv.KVStore, accountKey kv.Key, amount *big.Int) bool {
	balance := GetBaseTokensFullDecimals(v)(state, accountKey)
	if balance.Cmp(amount) < 0 {
		return false
	}
	setBaseTokensFullDecimals(v)(state, accountKey, new(big.Int).Sub(balance, amount))
	return true
}

// getFungibleTokens returns the fungible tokens owned by an account (base tokens extra decimals will be discarded)
func getFungibleTokens(v isc.SchemaVersion, state kv.KVStoreReader, accountKey kv.Key) *isc.Assets {
	ret := isc.NewEmptyAssets()
	ret.AddBaseTokens(getBaseTokens(v)(state, accountKey))
	NativeTokensMapR(state, accountKey).Iterate(func(idBytes []byte, val []byte) bool {
		ret.AddNativeTokens(
			isc.MustNativeTokenIDFromBytes(idBytes),
			new(big.Int).SetBytes(val),
		)
		return true
	})
	return ret
}

// GetAccountFungibleTokens returns all fungible tokens belonging to the agentID on the state
func GetAccountFungibleTokens(v isc.SchemaVersion, state kv.KVStoreReader, agentID isc.AgentID, chainID isc.ChainID) *isc.Assets {
	return getFungibleTokens(v, state, AccountKey(agentID, chainID))
}

func GetTotalL2FungibleTokens(v isc.SchemaVersion, state kv.KVStoreReader) *isc.Assets {
	return getFungibleTokens(v, state, L2TotalsAccount)
}

func getAccountBalanceDict(v isc.SchemaVersion, state kv.KVStoreReader, accountKey kv.Key) dict.Dict {
	return getFungibleTokens(v, state, accountKey).ToDict()
}
