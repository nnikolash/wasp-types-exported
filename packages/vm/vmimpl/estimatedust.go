package vmimpl

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/transaction"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
)

func (reqctx *requestContext) estimateRequiredStorageDeposit(par isc.RequestParameters) uint64 {
	par.AdjustToMinimumStorageDeposit = false

	hname := reqctx.CurrentContractHname()
	contractIdentity := isc.ContractIdentityFromHname(hname)
	if hname == evm.Contract.Hname() {
		contractIdentity = isc.ContractIdentityFromEVMAddress(common.Address{}) // use empty EVM address as STUB
	}
	out := transaction.BasicOutputFromPostData(
		reqctx.vm.task.AnchorOutput.AliasID.ToAddress(),
		contractIdentity,
		par,
	)
	return parameters.L1().Protocol.RentStructure.MinRent(out)
}
