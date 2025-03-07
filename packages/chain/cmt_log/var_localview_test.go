// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/chain/cmt_log"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testlogger"
)

func TestVarLocalView(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()
	j := cmt_log.NewVarLocalView(-1, func(ao *isc.AliasOutputWithID) {}, log)
	require.Nil(t, j.Value())
	tipAO, ok, _ := j.AliasOutputConfirmed(
		isc.NewAliasOutputWithID(
			&iotago.AliasOutput{
				StateMetadata: []byte{},
			},
			iotago.OutputID{},
		),
	)
	require.True(t, ok)
	require.NotNil(t, tipAO)
	require.NotNil(t, j.Value())
	require.Equal(t, tipAO, j.Value())
}
