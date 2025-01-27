// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coreroot

import (
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type MapStringToImmutableBytes struct {
	Proxy wasmtypes.Proxy
}

func (m MapStringToImmutableBytes) GetBytes(key string) wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(m.Proxy.Key(wasmtypes.StringToBytes(key)))
}

type ImmutableDeployContractParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableDeployContractParams() ImmutableDeployContractParams {
	return ImmutableDeployContractParams{Proxy: wasmlib.NewParamsProxy()}
}

// additional params for smart contract init function
func (s ImmutableDeployContractParams) InitParams() MapStringToImmutableBytes {
	return MapStringToImmutableBytes(s)
}

// The name of the contract to be deployed, used to calculate the contract's hname.
// The hname must be unique among all contract hnames in the chain.
func (s ImmutableDeployContractParams) Name() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.Proxy.Root(ParamName))
}

// hash of blob that has been previously stored in blob contract
func (s ImmutableDeployContractParams) ProgramHash() wasmtypes.ScImmutableHash {
	return wasmtypes.NewScImmutableHash(s.Proxy.Root(ParamProgramHash))
}

type MapStringToMutableBytes struct {
	Proxy wasmtypes.Proxy
}

func (m MapStringToMutableBytes) Clear() {
	m.Proxy.ClearMap()
}

func (m MapStringToMutableBytes) GetBytes(key string) wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(m.Proxy.Key(wasmtypes.StringToBytes(key)))
}

type MutableDeployContractParams struct {
	Proxy wasmtypes.Proxy
}

// additional params for smart contract init function
func (s MutableDeployContractParams) InitParams() MapStringToMutableBytes {
	return MapStringToMutableBytes(s)
}

// The name of the contract to be deployed, used to calculate the contract's hname.
// The hname must be unique among all contract hnames in the chain.
func (s MutableDeployContractParams) Name() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.Proxy.Root(ParamName))
}

// hash of blob that has been previously stored in blob contract
func (s MutableDeployContractParams) ProgramHash() wasmtypes.ScMutableHash {
	return wasmtypes.NewScMutableHash(s.Proxy.Root(ParamProgramHash))
}

type ImmutableGrantDeployPermissionParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableGrantDeployPermissionParams() ImmutableGrantDeployPermissionParams {
	return ImmutableGrantDeployPermissionParams{Proxy: wasmlib.NewParamsProxy()}
}

// agent to grant deploy permission to
func (s ImmutableGrantDeployPermissionParams) Deployer() wasmtypes.ScImmutableAgentID {
	return wasmtypes.NewScImmutableAgentID(s.Proxy.Root(ParamDeployer))
}

type MutableGrantDeployPermissionParams struct {
	Proxy wasmtypes.Proxy
}

// agent to grant deploy permission to
func (s MutableGrantDeployPermissionParams) Deployer() wasmtypes.ScMutableAgentID {
	return wasmtypes.NewScMutableAgentID(s.Proxy.Root(ParamDeployer))
}

type ImmutableRequireDeployPermissionsParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableRequireDeployPermissionsParams() ImmutableRequireDeployPermissionsParams {
	return ImmutableRequireDeployPermissionsParams{Proxy: wasmlib.NewParamsProxy()}
}

// turns permission check on or off
func (s ImmutableRequireDeployPermissionsParams) DeployPermissionsEnabled() wasmtypes.ScImmutableBool {
	return wasmtypes.NewScImmutableBool(s.Proxy.Root(ParamDeployPermissionsEnabled))
}

type MutableRequireDeployPermissionsParams struct {
	Proxy wasmtypes.Proxy
}

// turns permission check on or off
func (s MutableRequireDeployPermissionsParams) DeployPermissionsEnabled() wasmtypes.ScMutableBool {
	return wasmtypes.NewScMutableBool(s.Proxy.Root(ParamDeployPermissionsEnabled))
}

type ImmutableRevokeDeployPermissionParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableRevokeDeployPermissionParams() ImmutableRevokeDeployPermissionParams {
	return ImmutableRevokeDeployPermissionParams{Proxy: wasmlib.NewParamsProxy()}
}

// agent to revoke deploy permission for
func (s ImmutableRevokeDeployPermissionParams) Deployer() wasmtypes.ScImmutableAgentID {
	return wasmtypes.NewScImmutableAgentID(s.Proxy.Root(ParamDeployer))
}

type MutableRevokeDeployPermissionParams struct {
	Proxy wasmtypes.Proxy
}

// agent to revoke deploy permission for
func (s MutableRevokeDeployPermissionParams) Deployer() wasmtypes.ScMutableAgentID {
	return wasmtypes.NewScMutableAgentID(s.Proxy.Root(ParamDeployer))
}

type ImmutableFindContractParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableFindContractParams() ImmutableFindContractParams {
	return ImmutableFindContractParams{Proxy: wasmlib.NewParamsProxy()}
}

// The smart contract’s Hname
func (s ImmutableFindContractParams) Hname() wasmtypes.ScImmutableHname {
	return wasmtypes.NewScImmutableHname(s.Proxy.Root(ParamHname))
}

type MutableFindContractParams struct {
	Proxy wasmtypes.Proxy
}

// The smart contract’s Hname
func (s MutableFindContractParams) Hname() wasmtypes.ScMutableHname {
	return wasmtypes.NewScMutableHname(s.Proxy.Root(ParamHname))
}
