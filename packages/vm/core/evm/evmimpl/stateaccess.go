package evmimpl

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
)

type StateAccess struct {
	evmPartition kv.KVStoreReader
}

func NewStateAccess(store kv.KVStoreReader) *StateAccess {
	return &StateAccess{evmPartition: evm.ContractPartitionR(store)}
}

func (sa *StateAccess) Nonce(addr common.Address) uint64 {
	return Nonce(sa.evmPartition, addr)
}
