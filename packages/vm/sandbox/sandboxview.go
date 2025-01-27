// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package sandbox

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/vm/execution"
)

type sandboxView struct {
	SandboxBase
}

func NewSandboxView(ctx execution.WaspCallContext) isc.SandboxView {
	return &sandboxView{
		SandboxBase: SandboxBase{
			Ctx: ctx,
		},
	}
}
