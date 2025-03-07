// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

type inputStateMgrProposalConfirmed struct{}

func NewInputStateMgrProposalConfirmed() gpa.Input {
	return &inputStateMgrProposalConfirmed{}
}

func (inp *inputStateMgrProposalConfirmed) String() string {
	return "{cons.inputStateMgrProposalConfirmed}"
}
