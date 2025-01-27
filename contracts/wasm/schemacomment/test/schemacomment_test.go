// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/contracts/wasm/schemacomment/go/schemacomment"
	"github.com/nnikolash/wasp-types-exported/contracts/wasm/schemacomment/go/schemacommentimpl"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmsolo"
)

func TestDeploy(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, schemacomment.ScName, schemacommentimpl.OnDispatch)
	require.NoError(t, ctx.ContractExists(schemacomment.ScName))
}
