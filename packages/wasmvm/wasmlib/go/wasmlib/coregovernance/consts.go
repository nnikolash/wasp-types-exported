// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coregovernance

import "github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"

const (
	ScName        = "governance"
	ScDescription = "Governance contract"
	HScName       = wasmtypes.ScHname(0x17cf909f)
)

const (
	ParamAccessAPI     = "ia"
	ParamAccessOnly    = "i"
	ParamActions       = "n"
	ParamAddress       = "S"
	ParamCertificate   = "ic"
	ParamChainOwner    = "o"
	ParamFeePolicy     = "g"
	ParamGasLimits     = "l"
	ParamGasRatio      = "e"
	ParamMetadata      = "md"
	ParamPayoutAgentID = "s"
	ParamPubKey        = "ip"
	ParamPublicURL     = "x"
	ParamSetMinSD      = "ms"
)

const (
	ResultAccessNodeCandidates = "an"
	ResultAccessNodes          = "ac"
	ResultChainID              = "c"
	ResultChainOwnerID         = "o"
	ResultControllers          = "a"
	ResultFeePolicy            = "g"
	ResultGasLimits            = "l"
	ResultGasRatio             = "e"
	ResultGetMinSD             = "ms"
	ResultMetadata             = "md"
	ResultPayoutAgentID        = "s"
	ResultPublicURL            = "x"
	ResultStatus               = "m"
)

const (
	FuncAddAllowedStateControllerAddress    = "addAllowedStateControllerAddress"
	FuncAddCandidateNode                    = "addCandidateNode"
	FuncChangeAccessNodes                   = "changeAccessNodes"
	FuncClaimChainOwnership                 = "claimChainOwnership"
	FuncDelegateChainOwnership              = "delegateChainOwnership"
	FuncRemoveAllowedStateControllerAddress = "removeAllowedStateControllerAddress"
	FuncRevokeAccessNode                    = "revokeAccessNode"
	FuncRotateStateController               = "rotateStateController"
	FuncSetEVMGasRatio                      = "setEVMGasRatio"
	FuncSetFeePolicy                        = "setFeePolicy"
	FuncSetGasLimits                        = "setGasLimits"
	FuncSetMetadata                         = "setMetadata"
	FuncSetMinSD                            = "setMinSD"
	FuncSetPayoutAgentID                    = "setPayoutAgentID"
	FuncStartMaintenance                    = "startMaintenance"
	FuncStopMaintenance                     = "stopMaintenance"
	ViewGetAllowedStateControllerAddresses  = "getAllowedStateControllerAddresses"
	ViewGetChainInfo                        = "getChainInfo"
	ViewGetChainNodes                       = "getChainNodes"
	ViewGetChainOwner                       = "getChainOwner"
	ViewGetEVMGasRatio                      = "getEVMGasRatio"
	ViewGetFeePolicy                        = "getFeePolicy"
	ViewGetGasLimits                        = "getGasLimits"
	ViewGetMaintenanceStatus                = "getMaintenanceStatus"
	ViewGetMetadata                         = "getMetadata"
	ViewGetMinSD                            = "getMinSD"
	ViewGetPayoutAgentID                    = "getPayoutAgentID"
)

const (
	HFuncAddAllowedStateControllerAddress    = wasmtypes.ScHname(0x9469d567)
	HFuncAddCandidateNode                    = wasmtypes.ScHname(0xb745b382)
	HFuncChangeAccessNodes                   = wasmtypes.ScHname(0x7bca3700)
	HFuncClaimChainOwnership                 = wasmtypes.ScHname(0x03ff0fc0)
	HFuncDelegateChainOwnership              = wasmtypes.ScHname(0x93ecb6ad)
	HFuncRemoveAllowedStateControllerAddress = wasmtypes.ScHname(0x31f69447)
	HFuncRevokeAccessNode                    = wasmtypes.ScHname(0x5459512d)
	HFuncRotateStateController               = wasmtypes.ScHname(0x244d1038)
	HFuncSetEVMGasRatio                      = wasmtypes.ScHname(0xaae22338)
	HFuncSetFeePolicy                        = wasmtypes.ScHname(0x5b791c9f)
	HFuncSetGasLimits                        = wasmtypes.ScHname(0xd72fb355)
	HFuncSetMetadata                         = wasmtypes.ScHname(0x0eb3a798)
	HFuncSetMinSD                            = wasmtypes.ScHname(0x9cad5084)
	HFuncSetPayoutAgentID                    = wasmtypes.ScHname(0x2184ed1c)
	HFuncStartMaintenance                    = wasmtypes.ScHname(0x742f0521)
	HFuncStopMaintenance                     = wasmtypes.ScHname(0x4e017b6a)
	HViewGetAllowedStateControllerAddresses  = wasmtypes.ScHname(0xf3505183)
	HViewGetChainInfo                        = wasmtypes.ScHname(0x434477e2)
	HViewGetChainNodes                       = wasmtypes.ScHname(0xe1832289)
	HViewGetChainOwner                       = wasmtypes.ScHname(0x9b2ef0ac)
	HViewGetEVMGasRatio                      = wasmtypes.ScHname(0xb81c8c34)
	HViewGetFeePolicy                        = wasmtypes.ScHname(0xf8c89790)
	HViewGetGasLimits                        = wasmtypes.ScHname(0x3a493455)
	HViewGetMaintenanceStatus                = wasmtypes.ScHname(0x61fe5443)
	HViewGetMetadata                         = wasmtypes.ScHname(0x79ad1ac6)
	HViewGetMinSD                            = wasmtypes.ScHname(0x37f53a59)
	HViewGetPayoutAgentID                    = wasmtypes.ScHname(0x02aca9ad)
)
