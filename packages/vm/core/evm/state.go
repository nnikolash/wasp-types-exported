// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package evm

import (
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/subrealm"
)

// The evm core contract state stores two subrealms.
const (
	// KeyEmulatorState is the subrealm prefix for the data stored by the emulator (StateDB + BlockchainDB)
	KeyEmulatorState = "s"

	// KeyISCMagic is the subrealm prefix for the ISC magic contract
	KeyISCMagic = "m"
)

func ContractPartition(chainState kv.KVStore) kv.KVStore {
	return subrealm.New(chainState, kv.Key(Contract.Hname().Bytes()))
}

func ContractPartitionR(chainState kv.KVStoreReader) kv.KVStoreReader {
	return subrealm.NewReadOnly(chainState, kv.Key(Contract.Hname().Bytes()))
}

func EmulatorStateSubrealm(evmPartition kv.KVStore) kv.KVStore {
	return subrealm.New(evmPartition, KeyEmulatorState)
}

func EmulatorStateSubrealmR(evmPartition kv.KVStoreReader) kv.KVStoreReader {
	return subrealm.NewReadOnly(evmPartition, KeyEmulatorState)
}

func ISCMagicSubrealm(evmPartition kv.KVStore) kv.KVStore {
	return subrealm.New(evmPartition, KeyISCMagic)
}

func ISCMagicSubrealmR(evmPartition kv.KVStoreReader) kv.KVStoreReader {
	return subrealm.NewReadOnly(evmPartition, KeyISCMagic)
}
