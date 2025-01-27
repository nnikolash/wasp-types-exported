package m003

import (
	"github.com/iotaledger/hive.go/logger"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/emulator"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/iscmagic"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/migrations"
)

var UpdateEVMISCMagicFixed = migrations.Migration{
	Contract: evm.Contract,
	Apply: func(state kv.KVStore, log *logger.Logger) error {
		log.Infof("m003 UpdateEVMISCMagicFixed started")
		emulatorState := evm.EmulatorStateSubrealm(state)
		stateDBSubrealm := emulator.StateDBSubrealm(emulatorState)
		emulator.SetCode(stateDBSubrealm, iscmagic.ERC721NFTsAddress, iscmagic.ERC721NFTsRuntimeBytecode)
		log.Infof("m003 UpdateEVMISCMagicFixed finished")
		return nil
	},
}
