// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/documentation/tutorial-examples/go/solotutorial"
	"github.com/nnikolash/wasp-types-exported/documentation/tutorial-examples/go/solotutorialimpl"
	"github.com/nnikolash/wasp-types-exported/packages/wasmvm/wasmsolo"
)

func TestDeploy(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, solotutorial.ScName, solotutorialimpl.OnDispatch)
	require.NoError(t, ctx.ContractExists(solotutorial.ScName))
}
