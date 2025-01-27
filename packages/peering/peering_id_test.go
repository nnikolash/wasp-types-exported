// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package peering_test

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/peering"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestPeeringIDSerialization(t *testing.T) {
	peeringID := peering.RandomPeeringID()

	rwutil.ReadWriteTest(t, &peeringID, new(peering.PeeringID))
}
