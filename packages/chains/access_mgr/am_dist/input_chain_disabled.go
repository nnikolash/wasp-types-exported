package am_dist

import (
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
)

type inputChainDisabled struct {
	chainID isc.ChainID
}

func NewInputChainDisabled(chainID isc.ChainID) gpa.Input {
	return &inputChainDisabled{chainID: chainID}
}
