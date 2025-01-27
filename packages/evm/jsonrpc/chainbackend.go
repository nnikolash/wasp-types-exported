// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package jsonrpc

import (
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/trie"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

// ChainBackend provides access to the underlying ISC chain.
type ChainBackend interface {
	EVMSendTransaction(tx *types.Transaction) error
	EVMCall(aliasOutput *isc.AliasOutputWithID, callMsg ethereum.CallMsg) ([]byte, error)
	EVMEstimateGas(aliasOutput *isc.AliasOutputWithID, callMsg ethereum.CallMsg) (uint64, error)
	EVMTrace(aliasOutput *isc.AliasOutputWithID, blockTime time.Time, iscRequestsInBlock []isc.Request, txIndex *uint64, blockNumber *uint64, tracer *tracers.Tracer) error
	FeePolicy(blockIndex uint32) (*gas.FeePolicy, error)
	ISCChainID() *isc.ChainID
	ISCCallView(chainState state.State, scName string, funName string, args dict.Dict) (dict.Dict, error)
	ISCLatestAliasOutput() (*isc.AliasOutputWithID, error)
	ISCLatestState() (state.State, error)
	ISCStateByBlockIndex(blockIndex uint32) (state.State, error)
	ISCStateByTrieRoot(trieRoot trie.Hash) (state.State, error)
	BaseToken() *parameters.BaseToken
	TakeSnapshot() (int, error)
	RevertToSnapshot(int) error
}
