// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package generator

import "github.com/nnikolash/wasp-types-exported/tools/schema/model"

var commonTemplates = model.StringMap{
	// *******************************
	"else": `
`,
	// *******************************
	"newline": `

`,
	// *******************************
	"copyrightMessage": `
$#emit initGlobals
$#if copyright userCopyrightMessage
`,
	// *******************************
	"userCopyrightMessage": `
$copyright
`,
	// *******************************
	"warning": `
// Code generated by schema tool; DO NOT EDIT.

`,
	// *******************************
	"../README.md Lib": `
## $package

Interface library for: $scDesc
`,
	// *******************************
	"../README.md Impl": `
## $package$+impl

Implementation library for: $scDesc
`,
	// *******************************
	"../README.md Wasm": `
## $package$+wasm

Wasm VM host stub for: $scDesc
`,
	// *******************************
	"test.go": `
$#emit copyright
package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"$module/go/$package"
	"$module/go/$package$+impl"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmsolo"
)

func TestDeploy(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, $package.ScName, $package$+impl.OnDispatch)
	require.NoError(t, ctx.ContractExists($package.ScName))
}
`,
	// *******************************
	"setupInitFunc": `
$#set initFunc $nil
$#if init setInitFunc
`,
	// *******************************
	"setInitFunc": `
$#set initFunc Init
`,
	// *******************************
	"alignCalculate": `
$#set align $nil
$#if events align1space
$#if param align1space
$#if result align2spaces
$#set salign $align
$#set align $nil
$#if result align1space
$#set falign $nil
$#if salign alignSetFunc
`,
	// *******************************
	"align1space": `
$#set align $space
`,
	// *******************************
	"align2spaces": `
$#set align $space$space
`,
	// *******************************
	"alignSetFunc": `
$#set falign $salign$space
`,
}
