// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc_test

import (
	"crypto/rand"
	mathrand "math/rand"
	"testing"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestAliasOutputWithIDSerialization(t *testing.T) {
	output := iotago.AliasOutput{
		Amount:        mathrand.Uint64(),
		StateIndex:    mathrand.Uint32(),
		StateMetadata: []byte{},
	}
	rand.Read(output.AliasID[:])
	outputID := iotago.OutputID{}
	rand.Read(outputID[:])
	aliasOutputWithID := isc.NewAliasOutputWithID(&output, outputID)
	rwutil.ReadWriteTest(t, aliasOutputWithID, new(isc.AliasOutputWithID))
	rwutil.BytesTest(t, aliasOutputWithID, isc.AliasOutputWithIDFromBytes)
}
