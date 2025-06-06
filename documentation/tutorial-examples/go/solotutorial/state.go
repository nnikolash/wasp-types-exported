// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package solotutorial

import (
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type ImmutableSoloTutorialState struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableSoloTutorialState() ImmutableSoloTutorialState {
	return ImmutableSoloTutorialState{Proxy: wasmlib.NewStateProxy()}
}

func (s ImmutableSoloTutorialState) Str() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.Proxy.Root(StateStr))
}

type MutableSoloTutorialState struct {
	Proxy wasmtypes.Proxy
}

func NewMutableSoloTutorialState() MutableSoloTutorialState {
	return MutableSoloTutorialState{Proxy: wasmlib.NewStateProxy()}
}

func (s MutableSoloTutorialState) AsImmutable() ImmutableSoloTutorialState {
	return ImmutableSoloTutorialState(s)
}

func (s MutableSoloTutorialState) Str() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.Proxy.Root(StateStr))
}
