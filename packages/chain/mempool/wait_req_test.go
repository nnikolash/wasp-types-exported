// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package mempool_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/chain/mempool"
	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

func TestWaitReq(t *testing.T) {
	kp := cryptolib.NewKeyPair()

	ctxA := context.Background()
	ctxM := context.Background()
	req0 := isc.NewOffLedgerRequest(isc.RandomChainID(), isc.Hn("foo"), isc.Hn("bar"), nil, 0, gas.LimitsDefault.MaxGasPerRequest).Sign(kp)
	req1 := isc.NewOffLedgerRequest(isc.RandomChainID(), isc.Hn("foo"), isc.Hn("bar"), nil, 1, gas.LimitsDefault.MaxGasPerRequest).Sign(kp)
	req2 := isc.NewOffLedgerRequest(isc.RandomChainID(), isc.Hn("foo"), isc.Hn("bar"), nil, 2, gas.LimitsDefault.MaxGasPerRequest).Sign(kp)
	ref0 := isc.RequestRefFromRequest(req0)
	ref1 := isc.RequestRefFromRequest(req1)
	ref2 := isc.RequestRefFromRequest(req2)

	wr := mempool.NewWaitReq(3)
	var recvAny isc.Request
	recvMany := []isc.Request{}
	wr.WaitAny(ctxA, func(req isc.Request) {
		require.Nil(t, recvAny)
		recvAny = req
	})
	wr.WaitMany(ctxM, []*isc.RequestRef{ref0, ref1, ref2}, func(req isc.Request) {
		recvMany = append(recvMany, req)
	})
	wr.MarkAvailable(req0)
	wr.MarkAvailable(req1)
	wr.MarkAvailable(req2)
	require.NotNil(t, recvAny)
	require.Len(t, recvMany, 3)
	require.Contains(t, recvMany, req0)
	require.Contains(t, recvMany, req1)
	require.Contains(t, recvMany, req2)
}
