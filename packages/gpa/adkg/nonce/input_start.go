// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package nonce

import (
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

type inputStart struct{}

func NewInputStart() gpa.Input {
	return &inputStart{}
}
