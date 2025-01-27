// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package am_dist

import (
	"math"
	"math/rand"
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestMsgAccessSerialization(t *testing.T) {
	msg := &msgAccess{
		gpa.BasicMessage{},
		rand.Intn(math.MaxUint32 + 1),
		rand.Intn(math.MaxUint32 + 1),
		[]isc.ChainID{isc.RandomChainID(), isc.RandomChainID()},
		[]isc.ChainID{isc.RandomChainID(), isc.RandomChainID()},
	}

	rwutil.ReadWriteTest(t, msg, new(msgAccess))
}
