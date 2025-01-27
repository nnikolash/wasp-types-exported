// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package blocklog

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/kv/subrealm"
)

type StateAccess struct {
	state kv.KVStoreReader
}

func NewStateAccess(store kv.KVStoreReader) *StateAccess {
	state := subrealm.NewReadOnly(store, kv.Key(Contract.Hname().Bytes()))
	return &StateAccess{state: state}
}

func (sa *StateAccess) BlockInfo(blockIndex uint32) (*BlockInfo, bool) {
	return GetBlockInfo(sa.state, blockIndex)
}

func (sa *StateAccess) GetSmartContractEvents(contractID isc.Hname, fromBlock, toBlock uint32) dict.Dict {
	events := getSmartContractEventsInternal(sa.state, contractID, fromBlock, toBlock)
	return eventsToDict(events)
}
