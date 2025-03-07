package root

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc/coreutil"
)

var Contract = coreutil.NewContract(coreutil.CoreContractRoot)

var (
	// Funcs
	FuncDeployContract           = coreutil.Func("deployContract")
	FuncGrantDeployPermission    = coreutil.Func("grantDeployPermission")
	FuncRevokeDeployPermission   = coreutil.Func("revokeDeployPermission")
	FuncRequireDeployPermissions = coreutil.Func("requireDeployPermissions")

	// Views
	ViewFindContract       = coreutil.ViewFunc("findContract")
	ViewGetContractRecords = coreutil.ViewFunc("getContractRecords")
)

// state variables
const (
	VarSchemaVersion            = "v" // covered in: TestDeployNativeContract
	VarContractRegistry         = "r" // covered in: TestDeployNativeContract
	VarDeployPermissionsEnabled = "a" // covered in: TestDeployNativeContract
	VarDeployPermissions        = "p" // covered in: TestDeployNativeContract
)

// request parameters
const (
	ParamDeployer                 = "dp"
	ParamHname                    = "hn"
	ParamName                     = "nm"
	ParamProgramHash              = "ph"
	ParamContractRecData          = "dt"
	ParamContractFound            = "cf"
	ParamDeployPermissionsEnabled = "de"
)
