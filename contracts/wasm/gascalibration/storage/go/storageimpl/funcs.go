// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package storageimpl

import "github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmlib/go/wasmlib"

func funcF(_ wasmlib.ScFuncContext, f *FContext) {
	v := f.State.V()
	n := f.Params.N().Value()
	for i := uint32(0); i < n; i++ {
		v.AppendUint32().SetValue(i)
	}
}
