package rootimpl

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/collections"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/errors/coreerrors"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/root"
)

// isAuthorizedToDeploy checks if caller is authorized to deploy smart contract
func isAuthorizedToDeploy(ctx isc.Sandbox) bool {
	permissionsEnabled, err := codec.DecodeBool(ctx.State().Get(root.VarDeployPermissionsEnabled))
	if err != nil {
		return false
	}
	if !permissionsEnabled {
		return true
	}

	caller := ctx.Caller()
	if caller.Equals(ctx.ChainOwnerID()) {
		// chain owner is always authorized
		return true
	}
	if ctx.ChainID().IsSameChain(caller) {
		// smart contract from the same chain is always authorized
		return true
	}

	return collections.NewMap(ctx.State(), root.VarDeployPermissions).HasAt(caller.Bytes())
}

var errContractAlreadyExists = coreerrors.Register("contract with hname %08x already exists")

func storeContractRecord(state kv.KVStore, rec *root.ContractRecord) {
	hname := isc.Hn(rec.Name)
	// storing contract record in the registry
	contractRegistry := root.GetContractRegistry(state)
	if contractRegistry.HasAt(hname.Bytes()) {
		panic(errContractAlreadyExists.Create(hname))
	}
	contractRegistry.SetAt(hname.Bytes(), rec.Bytes())
}
