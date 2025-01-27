package evmimpl

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/nnikolash/wasp-types-exported/packages/evm/solidity"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/emulator"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/iscmagic"
)

func Nonce(evmPartition kv.KVStoreReader, addr common.Address) uint64 {
	emuState := evm.EmulatorStateSubrealmR(evmPartition)
	stateDBStore := emulator.StateDBSubrealmR(emuState)
	return emulator.GetNonce(stateDBStore, addr)
}

func registerERC721NFTCollectionByNFTId(evmState kv.KVStore, nft *isc.NFT) {
	metadata, err := isc.IRC27NFTMetadataFromBytes(nft.Metadata)
	if err != nil {
		panic(errEVMCanNotDecodeERC27Metadata)
	}

	addr := iscmagic.ERC721NFTCollectionAddress(nft.ID)
	state := emulator.NewStateDBFromKVStore(evm.EmulatorStateSubrealm(evmState))

	if state.Exist(addr) {
		panic(errEVMAccountAlreadyExists)
	}

	state.CreateAccount(addr)
	state.SetCode(addr, iscmagic.ERC721NFTCollectionRuntimeBytecode)
	// see ERC721NFTCollection_storage.json
	state.SetState(addr, solidity.StorageSlot(2), solidity.StorageEncodeBytes32(nft.ID[:]))
	for k, v := range solidity.StorageEncodeString(3, metadata.Name) {
		state.SetState(addr, k, v)
	}

	addToPrivileged(evmState, addr)
}
