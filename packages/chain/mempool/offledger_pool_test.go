package mempool

import (
	"math/big"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/testutil"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testkey"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testlogger"
)

func TestOffledgerMempoolAccountNonce(t *testing.T) {
	waitReq := NewWaitReq(waitRequestCleanupEvery)
	pool := NewOffledgerPool(100, waitReq, func(int) {}, func(time.Duration) {}, testlogger.NewSilentLogger("", true))

	// generate a bunch of requests for the same account
	kp, addr := testkey.GenKeyAddr()
	agentID := isc.NewAgentID(addr)

	req0 := testutil.DummyOffledgerRequestForAccount(isc.RandomChainID(), 0, kp)
	req1 := testutil.DummyOffledgerRequestForAccount(isc.RandomChainID(), 1, kp)
	req2 := testutil.DummyOffledgerRequestForAccount(isc.RandomChainID(), 2, kp)
	req2new := testutil.DummyOffledgerRequestForAccount(isc.RandomChainID(), 2, kp)
	pool.Add(req0)
	pool.Add(req1)
	pool.Add(req1) // try to add the same request many times
	pool.Add(req2)
	pool.Add(req1)
	require.EqualValues(t, 3, pool.refLUT.Size())
	require.EqualValues(t, 1, pool.reqsByAcountOrdered.Size())
	reqsInPoolForAccount, _ := pool.reqsByAcountOrdered.Get(agentID.String())
	require.Len(t, reqsInPoolForAccount, 3)
	pool.Add(req2new)
	pool.Add(req2new)
	require.EqualValues(t, 4, pool.refLUT.Size())
	require.EqualValues(t, 1, pool.reqsByAcountOrdered.Size())
	reqsInPoolForAccount, _ = pool.reqsByAcountOrdered.Get(agentID.String())
	require.Len(t, reqsInPoolForAccount, 4)

	// try to remove everything during iteration
	pool.Iterate(func(account string, entries []*OrderedPoolEntry) {
		for _, e := range entries {
			pool.Remove(e.req)
		}
	})
	require.EqualValues(t, 0, pool.refLUT.Size())
	require.EqualValues(t, 0, pool.reqsByAcountOrdered.Size())
}

func TestOffledgerMempoolLimit(t *testing.T) {
	waitReq := NewWaitReq(waitRequestCleanupEvery)
	poolSizeLimit := 3
	pool := NewOffledgerPool(poolSizeLimit, waitReq, func(int) {}, func(time.Duration) {}, testlogger.NewSilentLogger("", true))

	// create requests with different gas prices
	req0 := testutil.DummyEVMRequest(isc.RandomChainID(), big.NewInt(1))
	req1 := testutil.DummyEVMRequest(isc.RandomChainID(), big.NewInt(2))
	req2 := testutil.DummyEVMRequest(isc.RandomChainID(), big.NewInt(3))
	pool.Add(req0)
	pool.Add(req1)
	pool.Add(req2)

	assertPoolSize := func() {
		require.EqualValues(t, 3, pool.refLUT.Size())
		require.Len(t, pool.orderedByGasPrice, 3)
		require.EqualValues(t, 3, pool.reqsByAcountOrdered.Size())
	}
	contains := func(reqs ...isc.OffLedgerRequest) {
		for _, req := range reqs {
			lo.ContainsBy(pool.orderedByGasPrice, func(e *OrderedPoolEntry) bool {
				return e.req.ID().Equals(req.ID())
			})
		}
	}

	assertPoolSize()

	// add a request with high
	req3 := testutil.DummyEVMRequest(isc.RandomChainID(), big.NewInt(3))
	pool.Add(req3)
	assertPoolSize()
	contains(req1, req2, req3) // assert req3 was added and req0 was removed

	req4 := testutil.DummyEVMRequest(isc.RandomChainID(), big.NewInt(1))
	pool.Add(req4)
	assertPoolSize()
	contains(req1, req2, req3) // assert req4 is not added

	req5 := testutil.DummyEVMRequest(isc.RandomChainID(), big.NewInt(4))
	pool.Add(req5)
	assertPoolSize()

	contains(req2, req3, req5) // assert req5 was added and req1 was removed
}
