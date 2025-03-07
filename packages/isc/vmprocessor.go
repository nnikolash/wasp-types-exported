// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc

import (
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
)

// package present processor interface. It must be implemented by VM

// VMProcessor is an interface to the VM processor instance.
type VMProcessor interface {
	GetEntryPoint(code Hname) (VMProcessorEntryPoint, bool)
	Entrypoints() map[Hname]ProcessorEntryPoint
}

type ProcessorEntryPoint interface {
	VMProcessorEntryPoint
	Name() string
	Hname() Hname
}

// VMProcessorEntryPoint is an abstract interface by which VM is called by passing
// the Sandbox interface
type VMProcessorEntryPoint interface {
	Call(ctx interface{}) dict.Dict
	IsView() bool
}
