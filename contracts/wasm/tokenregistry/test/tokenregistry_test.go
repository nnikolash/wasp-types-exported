// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/contracts/wasm/tokenregistry/go/tokenregistry"
	"github.com/nnikolash/wasp-types-exported/contracts/wasm/tokenregistry/go/tokenregistryimpl"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmsolo"
)

func setupTest(t *testing.T) *wasmsolo.SoloContext {
	return wasmsolo.NewSoloContext(t, tokenregistry.ScName, tokenregistryimpl.OnDispatch)
}

func TestDeploy(t *testing.T) {
	ctx := setupTest(t)
	require.NoError(t, ctx.ContractExists(tokenregistry.ScName))
}
