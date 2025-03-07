// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coregovernance

import (
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type ImmutableAddAllowedStateControllerAddressParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableAddAllowedStateControllerAddressParams() ImmutableAddAllowedStateControllerAddressParams {
	return ImmutableAddAllowedStateControllerAddressParams{Proxy: wasmlib.NewParamsProxy()}
}

// state controller address
func (s ImmutableAddAllowedStateControllerAddressParams) Address() wasmtypes.ScImmutableAddress {
	return wasmtypes.NewScImmutableAddress(s.Proxy.Root(ParamAddress))
}

type MutableAddAllowedStateControllerAddressParams struct {
	Proxy wasmtypes.Proxy
}

// state controller address
func (s MutableAddAllowedStateControllerAddressParams) Address() wasmtypes.ScMutableAddress {
	return wasmtypes.NewScMutableAddress(s.Proxy.Root(ParamAddress))
}

type ImmutableAddCandidateNodeParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableAddCandidateNodeParams() ImmutableAddCandidateNodeParams {
	return ImmutableAddCandidateNodeParams{Proxy: wasmlib.NewParamsProxy()}
}

// API base URL for the node, default empty
func (s ImmutableAddCandidateNodeParams) AccessAPI() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.Proxy.Root(ParamAccessAPI))
}

// whether node is just an access node, default false
func (s ImmutableAddCandidateNodeParams) AccessOnly() wasmtypes.ScImmutableBool {
	return wasmtypes.NewScImmutableBool(s.Proxy.Root(ParamAccessOnly))
}

// signed binary containing both the node public key and their L1 address
func (s ImmutableAddCandidateNodeParams) Certificate() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamCertificate))
}

// public key of the node to be added
func (s ImmutableAddCandidateNodeParams) PubKey() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamPubKey))
}

type MutableAddCandidateNodeParams struct {
	Proxy wasmtypes.Proxy
}

// API base URL for the node, default empty
func (s MutableAddCandidateNodeParams) AccessAPI() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.Proxy.Root(ParamAccessAPI))
}

// whether node is just an access node, default false
func (s MutableAddCandidateNodeParams) AccessOnly() wasmtypes.ScMutableBool {
	return wasmtypes.NewScMutableBool(s.Proxy.Root(ParamAccessOnly))
}

// signed binary containing both the node public key and their L1 address
func (s MutableAddCandidateNodeParams) Certificate() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamCertificate))
}

// public key of the node to be added
func (s MutableAddCandidateNodeParams) PubKey() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamPubKey))
}

type MapBytesToImmutableUint8 struct {
	Proxy wasmtypes.Proxy
}

func (m MapBytesToImmutableUint8) GetUint8(key []byte) wasmtypes.ScImmutableUint8 {
	return wasmtypes.NewScImmutableUint8(m.Proxy.Key(wasmtypes.BytesToBytes(key)))
}

type ImmutableChangeAccessNodesParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableChangeAccessNodesParams() ImmutableChangeAccessNodesParams {
	return ImmutableChangeAccessNodesParams{Proxy: wasmlib.NewParamsProxy()}
}

// map of actions per pubkey
// 0: Remove the access node from the access nodes list.
// 1: Accept a candidate node and add it to the list of access nodes.
// 2: Drop an access node from the access node and candidate lists.
func (s ImmutableChangeAccessNodesParams) Actions() MapBytesToImmutableUint8 {
	return MapBytesToImmutableUint8{Proxy: s.Proxy.Root(ParamActions)}
}

type MapBytesToMutableUint8 struct {
	Proxy wasmtypes.Proxy
}

func (m MapBytesToMutableUint8) Clear() {
	m.Proxy.ClearMap()
}

func (m MapBytesToMutableUint8) GetUint8(key []byte) wasmtypes.ScMutableUint8 {
	return wasmtypes.NewScMutableUint8(m.Proxy.Key(wasmtypes.BytesToBytes(key)))
}

type MutableChangeAccessNodesParams struct {
	Proxy wasmtypes.Proxy
}

// map of actions per pubkey
// 0: Remove the access node from the access nodes list.
// 1: Accept a candidate node and add it to the list of access nodes.
// 2: Drop an access node from the access node and candidate lists.
func (s MutableChangeAccessNodesParams) Actions() MapBytesToMutableUint8 {
	return MapBytesToMutableUint8{Proxy: s.Proxy.Root(ParamActions)}
}

type ImmutableDelegateChainOwnershipParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableDelegateChainOwnershipParams() ImmutableDelegateChainOwnershipParams {
	return ImmutableDelegateChainOwnershipParams{Proxy: wasmlib.NewParamsProxy()}
}

// next chain owner's agent ID
func (s ImmutableDelegateChainOwnershipParams) ChainOwner() wasmtypes.ScImmutableAgentID {
	return wasmtypes.NewScImmutableAgentID(s.Proxy.Root(ParamChainOwner))
}

type MutableDelegateChainOwnershipParams struct {
	Proxy wasmtypes.Proxy
}

// next chain owner's agent ID
func (s MutableDelegateChainOwnershipParams) ChainOwner() wasmtypes.ScMutableAgentID {
	return wasmtypes.NewScMutableAgentID(s.Proxy.Root(ParamChainOwner))
}

type ImmutableRemoveAllowedStateControllerAddressParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableRemoveAllowedStateControllerAddressParams() ImmutableRemoveAllowedStateControllerAddressParams {
	return ImmutableRemoveAllowedStateControllerAddressParams{Proxy: wasmlib.NewParamsProxy()}
}

// state controller address
func (s ImmutableRemoveAllowedStateControllerAddressParams) Address() wasmtypes.ScImmutableAddress {
	return wasmtypes.NewScImmutableAddress(s.Proxy.Root(ParamAddress))
}

type MutableRemoveAllowedStateControllerAddressParams struct {
	Proxy wasmtypes.Proxy
}

// state controller address
func (s MutableRemoveAllowedStateControllerAddressParams) Address() wasmtypes.ScMutableAddress {
	return wasmtypes.NewScMutableAddress(s.Proxy.Root(ParamAddress))
}

type ImmutableRevokeAccessNodeParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableRevokeAccessNodeParams() ImmutableRevokeAccessNodeParams {
	return ImmutableRevokeAccessNodeParams{Proxy: wasmlib.NewParamsProxy()}
}

// certificate of the node to be removed
func (s ImmutableRevokeAccessNodeParams) Certificate() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamCertificate))
}

// public key of the node to be removed
func (s ImmutableRevokeAccessNodeParams) PubKey() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamPubKey))
}

type MutableRevokeAccessNodeParams struct {
	Proxy wasmtypes.Proxy
}

// certificate of the node to be removed
func (s MutableRevokeAccessNodeParams) Certificate() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamCertificate))
}

// public key of the node to be removed
func (s MutableRevokeAccessNodeParams) PubKey() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamPubKey))
}

type ImmutableRotateStateControllerParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableRotateStateControllerParams() ImmutableRotateStateControllerParams {
	return ImmutableRotateStateControllerParams{Proxy: wasmlib.NewParamsProxy()}
}

// state controller address
func (s ImmutableRotateStateControllerParams) Address() wasmtypes.ScImmutableAddress {
	return wasmtypes.NewScImmutableAddress(s.Proxy.Root(ParamAddress))
}

type MutableRotateStateControllerParams struct {
	Proxy wasmtypes.Proxy
}

// state controller address
func (s MutableRotateStateControllerParams) Address() wasmtypes.ScMutableAddress {
	return wasmtypes.NewScMutableAddress(s.Proxy.Root(ParamAddress))
}

type ImmutableSetEVMGasRatioParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableSetEVMGasRatioParams() ImmutableSetEVMGasRatioParams {
	return ImmutableSetEVMGasRatioParams{Proxy: wasmlib.NewParamsProxy()}
}

// serialized gas ratio
func (s ImmutableSetEVMGasRatioParams) GasRatio() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamGasRatio))
}

type MutableSetEVMGasRatioParams struct {
	Proxy wasmtypes.Proxy
}

// serialized gas ratio
func (s MutableSetEVMGasRatioParams) GasRatio() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamGasRatio))
}

type ImmutableSetFeePolicyParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableSetFeePolicyParams() ImmutableSetFeePolicyParams {
	return ImmutableSetFeePolicyParams{Proxy: wasmlib.NewParamsProxy()}
}

// serialized fee policy
func (s ImmutableSetFeePolicyParams) FeePolicy() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamFeePolicy))
}

type MutableSetFeePolicyParams struct {
	Proxy wasmtypes.Proxy
}

// serialized fee policy
func (s MutableSetFeePolicyParams) FeePolicy() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamFeePolicy))
}

type ImmutableSetGasLimitsParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableSetGasLimitsParams() ImmutableSetGasLimitsParams {
	return ImmutableSetGasLimitsParams{Proxy: wasmlib.NewParamsProxy()}
}

// serialized gas limits
func (s ImmutableSetGasLimitsParams) GasLimits() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamGasLimits))
}

type MutableSetGasLimitsParams struct {
	Proxy wasmtypes.Proxy
}

// serialized gas limits
func (s MutableSetGasLimitsParams) GasLimits() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamGasLimits))
}

type ImmutableSetMetadataParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableSetMetadataParams() ImmutableSetMetadataParams {
	return ImmutableSetMetadataParams{Proxy: wasmlib.NewParamsProxy()}
}

// the public evm json rpc url
func (s ImmutableSetMetadataParams) Metadata() ImmutablePublicChainMetadata {
	return ImmutablePublicChainMetadata{Proxy: s.Proxy.Root(ParamMetadata)}
}

// the public url leading to the chain info, stored on the tangle
func (s ImmutableSetMetadataParams) PublicURL() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.Proxy.Root(ParamPublicURL))
}

type MutableSetMetadataParams struct {
	Proxy wasmtypes.Proxy
}

// the public evm json rpc url
func (s MutableSetMetadataParams) Metadata() MutablePublicChainMetadata {
	return MutablePublicChainMetadata{Proxy: s.Proxy.Root(ParamMetadata)}
}

// the public url leading to the chain info, stored on the tangle
func (s MutableSetMetadataParams) PublicURL() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.Proxy.Root(ParamPublicURL))
}

type ImmutableSetMinSDParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableSetMinSDParams() ImmutableSetMinSDParams {
	return ImmutableSetMinSDParams{Proxy: wasmlib.NewParamsProxy()}
}

func (s ImmutableSetMinSDParams) SetMinSD() wasmtypes.ScImmutableUint64 {
	return wasmtypes.NewScImmutableUint64(s.Proxy.Root(ParamSetMinSD))
}

type MutableSetMinSDParams struct {
	Proxy wasmtypes.Proxy
}

func (s MutableSetMinSDParams) SetMinSD() wasmtypes.ScMutableUint64 {
	return wasmtypes.NewScMutableUint64(s.Proxy.Root(ParamSetMinSD))
}

type ImmutableSetPayoutAgentIDParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableSetPayoutAgentIDParams() ImmutableSetPayoutAgentIDParams {
	return ImmutableSetPayoutAgentIDParams{Proxy: wasmlib.NewParamsProxy()}
}

// set payout AgentID
func (s ImmutableSetPayoutAgentIDParams) PayoutAgentID() wasmtypes.ScImmutableAgentID {
	return wasmtypes.NewScImmutableAgentID(s.Proxy.Root(ParamPayoutAgentID))
}

type MutableSetPayoutAgentIDParams struct {
	Proxy wasmtypes.Proxy
}

// set payout AgentID
func (s MutableSetPayoutAgentIDParams) PayoutAgentID() wasmtypes.ScMutableAgentID {
	return wasmtypes.NewScMutableAgentID(s.Proxy.Root(ParamPayoutAgentID))
}
