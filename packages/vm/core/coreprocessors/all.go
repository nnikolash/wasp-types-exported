package coreprocessors

import (
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blob"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blocklog"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/corecontracts"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/errors"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/evmimpl"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/governance"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/governance/governanceimpl"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/root"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/root/rootimpl"
	"github.com/nnikolash/wasp-types-exported/packages/vm/processors"
)

var All = map[hashing.HashValue]isc.VMProcessor{
	root.Contract.ProgramHash:       rootimpl.Processor,
	errors.Contract.ProgramHash:     errors.Processor,
	accounts.Contract.ProgramHash:   accounts.Processor,
	blob.Contract.ProgramHash:       blob.Processor,
	blocklog.Contract.ProgramHash:   blocklog.Processor,
	governance.Contract.ProgramHash: governanceimpl.Processor,
	evm.Contract.ProgramHash:        evmimpl.Processor,
}

func init() {
	if len(corecontracts.All) != len(All) {
		panic("static check: mismatch between corecontracts.All and coreprocessors.All")
	}
}

func NewConfigWithCoreContracts() *processors.Config {
	return processors.NewConfig().WithCoreContracts(All)
}
