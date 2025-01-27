// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package distsync

import "github.com/nnikolash/wasp-types-exported/packages/gpa"

type inputTimeTick struct{}

func NewInputTimeTick() gpa.Input {
	return &inputTimeTick{}
}
