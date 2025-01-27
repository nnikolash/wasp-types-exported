// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dss

import (
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

type inputStart struct{}

func NewInputStart() gpa.Input {
	return &inputStart{}
}

func (inp *inputStart) String() string {
	return "{chain.dss.inputStart}"
}
