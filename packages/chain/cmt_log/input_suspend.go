// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

// Sent by the Chain, when another committee has received control over this chain.
type inputSuspend struct{}

func NewInputSuspend() gpa.Input {
	return &inputSuspend{}
}

func (inp *inputSuspend) String() string {
	return "{cmtLog.inputSuspend}"
}
