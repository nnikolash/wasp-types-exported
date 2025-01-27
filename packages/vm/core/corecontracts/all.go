package corecontracts

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/isc/coreutil"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blob"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blocklog"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/errors"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/governance"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/root"
)

var All = map[isc.Hname]*coreutil.ContractInfo{
	root.Contract.Hname():       root.Contract,
	errors.Contract.Hname():     errors.Contract,
	accounts.Contract.Hname():   accounts.Contract,
	blob.Contract.Hname():       blob.Contract,
	blocklog.Contract.Hname():   blocklog.Contract,
	governance.Contract.Hname(): governance.Contract,
	evm.Contract.Hname():        evm.Contract,
}

func IsCoreHname(hname isc.Hname) bool {
	_, ok := All[hname]
	return ok
}
