package chainutil

import (
	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/vm/viewcontext"
)

// CallView executes a view call on the latest block of the chain
func CallView(
	chainState state.State,
	ch chain.ChainCore,
	contractHname,
	viewHname isc.Hname,
	params dict.Dict,
) (dict.Dict, error) {
	vctx, err := viewcontext.New(ch, chainState, false)
	if err != nil {
		return nil, err
	}
	return vctx.CallViewExternal(contractHname, viewHname, params)
}
