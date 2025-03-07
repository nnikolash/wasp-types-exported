package isc_test

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestVMErrorCodeSerialization(t *testing.T) {
	vmErrorCode := isc.VMErrorCode{
		ContractID: isc.Hname(1074),
		ID:         123,
	}
	rwutil.ReadWriteTest(t, &vmErrorCode, new(isc.VMErrorCode))
	rwutil.BytesTest(t, vmErrorCode, isc.VMErrorCodeFromBytes)
}
