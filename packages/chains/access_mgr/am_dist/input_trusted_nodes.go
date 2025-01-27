// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package am_dist

import (
	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

type inputTrustedNodes struct {
	trustedNodes []*cryptolib.PublicKey
}

var _ gpa.Input = &inputTrustedNodes{}

func NewInputTrustedNodes(trustedNodes []*cryptolib.PublicKey) gpa.Input {
	return &inputTrustedNodes{trustedNodes: trustedNodes}
}
