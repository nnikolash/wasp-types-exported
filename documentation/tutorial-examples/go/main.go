// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

//go:build wasm
// +build wasm

package main

import (
	"github.com/nnikolash/wasp-types-exported/documentation/tutorial-examples/go/solotutorialimpl"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmvmhost/go/wasmvmhost"
)

func main() {
}

func init() {
	wasmvmhost.ConnectWasmHost()
}

//export on_call
func onCall(index int32) {
	solotutorialimpl.OnDispatch(index)
}

//export on_load
func onLoad() {
	solotutorialimpl.OnDispatch(-1)
}
