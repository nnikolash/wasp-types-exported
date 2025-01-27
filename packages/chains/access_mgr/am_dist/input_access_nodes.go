// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package am_dist

import (
	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
)

type inputAccessNodes struct {
	chainID     isc.ChainID
	accessNodes []*cryptolib.PublicKey
}

var _ gpa.Input = &inputAccessNodes{}

func NewInputAccessNodes(chainID isc.ChainID, accessNodes []*cryptolib.PublicKey) gpa.Input {
	return &inputAccessNodes{chainID: chainID, accessNodes: accessNodes}
}
