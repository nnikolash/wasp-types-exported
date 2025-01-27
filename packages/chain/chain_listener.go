// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"github.com/nnikolash/wasp-types-exported/packages/chain/mempool"
	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/state"
)

// Implementation of this interface will receive events in the chain.
// Initial intention was to provide info to the published/WebSocket endpoint.
// All the function MUST NOT BLOCK.
type ChainListener interface {
	mempool.ChainListener
	AccessNodesUpdated(chainID isc.ChainID, accessNodes []*cryptolib.PublicKey)
	ServerNodesUpdated(chainID isc.ChainID, serverNodes []*cryptolib.PublicKey)
}

////////////////////////////////////////////////////////////////////////////////
// emptyChainListener

type emptyChainListener struct{}

var _ ChainListener = &emptyChainListener{}

func NewEmptyChainListener() ChainListener {
	return &emptyChainListener{}
}

func (ecl *emptyChainListener) BlockApplied(chainID isc.ChainID, block state.Block, latestState kv.KVStoreReader) {
}
func (ecl *emptyChainListener) AccessNodesUpdated(isc.ChainID, []*cryptolib.PublicKey) {}
func (ecl *emptyChainListener) ServerNodesUpdated(isc.ChainID, []*cryptolib.PublicKey) {}
