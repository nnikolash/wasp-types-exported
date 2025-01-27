// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package governanceimpl

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/governance"
)

// getChainInfo view returns general info about the chain: chain ID, chain owner ID, limits and default fees
func getChainInfo(ctx isc.SandboxView) dict.Dict {
	info := governance.MustGetChainInfo(ctx.StateR(), ctx.ChainID())
	ret := dict.New()
	ret.Set(governance.ParamChainID, codec.EncodeChainID(info.ChainID))
	ret.Set(governance.VarChainOwnerID, codec.EncodeAgentID(info.ChainOwnerID))
	ret.Set(governance.VarGasFeePolicyBytes, info.GasFeePolicy.Bytes())
	ret.Set(governance.VarGasLimitsBytes, info.GasLimits.Bytes())

	if info.PublicURL != "" {
		ret.Set(governance.VarPublicURL, codec.EncodeString(info.PublicURL))
	}

	ret.Set(governance.VarMetadata, info.Metadata.Bytes())

	return ret
}
