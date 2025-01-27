package chainutil

import (
	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/errors"
)

func ResolveError(ch chain.ChainCore, e *isc.UnresolvedVMError) (*isc.VMError, error) {
	s, err := ch.LatestState(chain.ActiveOrCommittedState)
	if err != nil {
		return nil, err
	}
	return errors.ResolveFromState(s, e)
}
