package vmimpl

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/subrealm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/execution"
)

func (reqctx *requestContext) chainStateWithGasBurn() kv.KVStore {
	return execution.NewKVStoreWithGasBurn(reqctx.uncommittedState, reqctx)
}

func (reqctx *requestContext) contractStateWithGasBurn() kv.KVStore {
	return subrealm.New(reqctx.chainStateWithGasBurn(), kv.Key(reqctx.CurrentContractHname().Bytes()))
}

func (reqctx *requestContext) ContractStateReaderWithGasBurn() kv.KVStoreReader {
	return subrealm.NewReadOnly(reqctx.chainStateWithGasBurn(), kv.Key(reqctx.CurrentContractHname().Bytes()))
}

func (reqctx *requestContext) SchemaVersion() isc.SchemaVersion {
	return reqctx.vm.schemaVersion
}
