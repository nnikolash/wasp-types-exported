package execution

import (
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

type kvStoreWithGasBurn struct {
	kv.KVStore
	gas GasContext
}

func NewKVStoreWithGasBurn(s kv.KVStore, gas GasContext) kv.KVStore {
	return &kvStoreWithGasBurn{
		KVStore: s,
		gas:     gas,
	}
}

func (s *kvStoreWithGasBurn) Get(name kv.Key) []byte {
	return getWithGasBurn(s.KVStore, name, s.gas)
}

func (s *kvStoreWithGasBurn) Set(name kv.Key, value []byte) {
	s.KVStore.Set(name, value)
	s.gas.GasBurn(gas.BurnCodeStorage1P, uint64(len(name)+len(value)))
}

type kvStoreReaderWithGasBurn struct {
	kv.KVStoreReader
	gas GasContext
}

func NewKVStoreReaderWithGasBurn(r kv.KVStoreReader, gas GasContext) kv.KVStoreReader {
	return &kvStoreReaderWithGasBurn{
		KVStoreReader: r,
		gas:           gas,
	}
}

func (s *kvStoreReaderWithGasBurn) Get(name kv.Key) []byte {
	return getWithGasBurn(s.KVStoreReader, name, s.gas)
}

func getWithGasBurn(r kv.KVStoreReader, name kv.Key, gasctx GasContext) []byte {
	v := r.Get(name)
	gasctx.GasBurn(gas.BurnCodeReadFromState1P, uint64(len(v)/100)+1) // minimum 1
	return v
}
