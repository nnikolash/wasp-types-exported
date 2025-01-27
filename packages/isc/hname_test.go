// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc_test

import (
	"math/rand"
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestHnameSerialize(t *testing.T) {
	hname := isc.Hname(rand.Uint32())
	rwutil.ReadWriteTest(t, &hname, new(isc.Hname))
	rwutil.BytesTest(t, hname, isc.HnameFromBytes)
	rwutil.StringTest(t, hname, isc.HnameFromString)
}
