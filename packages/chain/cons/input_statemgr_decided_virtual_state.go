// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"fmt"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/state"
)

type inputStateMgrDecidedVirtualState struct {
	chainState state.State
}

func NewInputStateMgrDecidedVirtualState(
	chainState state.State,
) gpa.Input {
	return &inputStateMgrDecidedVirtualState{chainState: chainState}
}

func (inp *inputStateMgrDecidedVirtualState) String() string {
	return fmt.Sprintf(
		"{cons.inputStateMgrDecidedVirtualState: blockIndex=%v, trieRoot=%v}",
		inp.chainState.BlockIndex(),
		inp.chainState.TrieRoot(),
	)
}
