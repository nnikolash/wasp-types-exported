// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package governance

import (
	"math/big"

	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/subrealm"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
)

type StateAccess struct {
	state kv.KVStoreReader
}

func NewStateAccess(store kv.KVStoreReader) *StateAccess {
	state := subrealm.NewReadOnly(store, kv.Key(Contract.Hname().Bytes()))
	return &StateAccess{state: state}
}

func (sa *StateAccess) MaintenanceStatus() bool {
	r := sa.state.Get(VarMaintenanceStatus)
	if r == nil {
		return false // chain is being initialized, governance has not been initialized yet
	}
	return codec.MustDecodeBool(r)
}

func (sa *StateAccess) AccessNodes() []*cryptolib.PublicKey {
	accessNodes := []*cryptolib.PublicKey{}
	AccessNodesMapR(sa.state).IterateKeys(func(pubKeyBytes []byte) bool {
		pubKey, err := cryptolib.PublicKeyFromBytes(pubKeyBytes)
		if err != nil {
			panic(err)
		}
		accessNodes = append(accessNodes, pubKey)
		return true
	})
	return accessNodes
}

func (sa *StateAccess) CandidateNodes() []*AccessNodeInfo {
	candidateNodes := []*AccessNodeInfo{}
	AccessNodeCandidatesMapR(sa.state).Iterate(func(pubKeyBytes, accessNodeInfoBytes []byte) bool {
		ani, err := AccessNodeInfoFromBytes(pubKeyBytes, accessNodeInfoBytes)
		if err != nil {
			panic(err)
		}
		candidateNodes = append(candidateNodes, ani)
		return true
	})
	return candidateNodes
}

func (sa *StateAccess) ChainInfo(chainID isc.ChainID) *isc.ChainInfo {
	return MustGetChainInfo(sa.state, chainID)
}

func (sa *StateAccess) ChainOwnerID() isc.AgentID {
	return mustGetChainOwnerID(sa.state)
}

func (sa *StateAccess) GetBlockKeepAmount() int32 {
	return GetBlockKeepAmount(sa.state)
}

func (sa *StateAccess) DefaultGasPrice() *big.Int {
	return MustGetGasFeePolicy(sa.state).DefaultGasPriceFullDecimals(parameters.L1().BaseToken.Decimals)
}
