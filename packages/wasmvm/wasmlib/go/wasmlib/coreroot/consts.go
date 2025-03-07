// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coreroot

import "github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"

const (
	ScName        = "root"
	ScDescription = "Root Contract"
	HScName       = wasmtypes.ScHname(0xcebf5908)
)

const (
	ParamDeployer                 = "dp"
	ParamDeployPermissionsEnabled = "de"
	ParamHname                    = "hn"
	ParamInitParams               = "this"
	ParamName                     = "nm"
	ParamProgramHash              = "ph"
)

const (
	ResultContractFound    = "cf"
	ResultContractRecData  = "dt"
	ResultContractRegistry = "r"
)

const (
	FuncDeployContract           = "deployContract"
	FuncGrantDeployPermission    = "grantDeployPermission"
	FuncRequireDeployPermissions = "requireDeployPermissions"
	FuncRevokeDeployPermission   = "revokeDeployPermission"
	ViewFindContract             = "findContract"
	ViewGetContractRecords       = "getContractRecords"
)

const (
	HFuncDeployContract           = wasmtypes.ScHname(0x28232c27)
	HFuncGrantDeployPermission    = wasmtypes.ScHname(0xf440263a)
	HFuncRequireDeployPermissions = wasmtypes.ScHname(0xefff8d83)
	HFuncRevokeDeployPermission   = wasmtypes.ScHname(0x850744f1)
	HViewFindContract             = wasmtypes.ScHname(0xc145ca00)
	HViewGetContractRecords       = wasmtypes.ScHname(0x078b3ef3)
)
