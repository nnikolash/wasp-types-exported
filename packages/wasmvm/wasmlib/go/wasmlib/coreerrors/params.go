// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coreerrors

import (
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type ImmutableRegisterErrorParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableRegisterErrorParams() ImmutableRegisterErrorParams {
	return ImmutableRegisterErrorParams{Proxy: wasmlib.NewParamsProxy()}
}

// error message template string
func (s ImmutableRegisterErrorParams) Template() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.Proxy.Root(ParamTemplate))
}

type MutableRegisterErrorParams struct {
	Proxy wasmtypes.Proxy
}

// error message template string
func (s MutableRegisterErrorParams) Template() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.Proxy.Root(ParamTemplate))
}

type ImmutableGetErrorMessageFormatParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableGetErrorMessageFormatParams() ImmutableGetErrorMessageFormatParams {
	return ImmutableGetErrorMessageFormatParams{Proxy: wasmlib.NewParamsProxy()}
}

// serialized error code
func (s ImmutableGetErrorMessageFormatParams) ErrorCode() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamErrorCode))
}

type MutableGetErrorMessageFormatParams struct {
	Proxy wasmtypes.Proxy
}

// serialized error code
func (s MutableGetErrorMessageFormatParams) ErrorCode() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamErrorCode))
}
