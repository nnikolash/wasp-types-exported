package blob

import (
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc/coreutil"
	"github.com/nnikolash/wasp-types-exported/packages/kv/collections"
)

var Contract = coreutil.NewContract(coreutil.CoreContractBlob)

var (
	FuncStoreBlob = coreutil.Func("storeBlob")

	ViewGetBlobInfo  = coreutil.ViewFunc("getBlobInfo")
	ViewGetBlobField = coreutil.ViewFunc("getBlobField")
)

// state variables
const (
	// variable names of standard blob's field
	// user-defined field must be different
	VarFieldProgramBinary      = "p"
	VarFieldVMType             = "v"
	VarFieldProgramDescription = "d"
)

// request parameters
const (
	ParamHash  = "hash"
	ParamField = "field"
	ParamBytes = "bytes"
)

// FieldValueKey returns key of the blob field value in the SC state.
func FieldValueKey(blobHash hashing.HashValue, fieldName string) []byte {
	return []byte(collections.MapElemKey(valuesMapName(blobHash), []byte(fieldName)))
}
