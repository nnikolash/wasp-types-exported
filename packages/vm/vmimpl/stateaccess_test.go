package vmimpl

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/buffered"
	"github.com/nnikolash/wasp-types-exported/packages/origin"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/vm"
)

func TestSetThenGet(t *testing.T) {
	db := mapdb.NewMapDB()
	cs := state.NewStoreWithUniqueWriteMutex(db)
	origin.InitChain(0, cs, nil, 0)
	latest, err := cs.LatestBlock()
	require.NoError(t, err)
	stateDraft, err := cs.NewStateDraft(time.Now(), latest.L1Commitment())
	require.NoError(t, err)

	hname := isc.Hn("test")

	reqctx := &requestContext{
		vm: &vmContext{
			task:       &vm.VMTask{},
			stateDraft: stateDraft,
		},
		uncommittedState: buffered.NewBufferedKVStore(stateDraft),
		callStack:        []*callContext{{contract: hname}},
	}
	s := reqctx.contractStateWithGasBurn()

	subpartitionedKey := kv.Key(hname.Bytes()) + "x"

	// contract sets variable x
	s.Set("x", []byte{42})
	require.Equal(t, map[kv.Key][]byte{subpartitionedKey: {42}}, reqctx.uncommittedState.Mutations().Sets)
	require.Equal(t, map[kv.Key]struct{}{}, reqctx.uncommittedState.Mutations().Dels)

	// contract gets variable x
	v := s.Get("x")
	require.Equal(t, []byte{42}, v)

	// mutation is in currentStateUpdate, prefixed by the contract id
	require.Equal(t, []byte{42}, reqctx.uncommittedState.Mutations().Sets[subpartitionedKey])

	// mutation is in the not committed to the virtual state yet
	v = stateDraft.Get(subpartitionedKey)
	require.Nil(t, v)

	// contract deletes variable x
	s.Del("x")
	require.Equal(t, map[kv.Key][]byte{}, reqctx.uncommittedState.Mutations().Sets)
	require.Equal(t, map[kv.Key]struct{}{subpartitionedKey: {}}, reqctx.uncommittedState.Mutations().Dels)

	// contract sees variable x does not exist
	v = s.Get("x")
	require.Nil(t, v)

	// contract makes several writes to same variable, gets the latest value
	s.Set("x", []byte{2 * 42})
	require.Equal(t, map[kv.Key][]byte{subpartitionedKey: {2 * 42}}, reqctx.uncommittedState.Mutations().Sets)
	require.Equal(t, map[kv.Key]struct{}{}, reqctx.uncommittedState.Mutations().Dels)

	s.Set("x", []byte{3 * 42})
	require.Equal(t, map[kv.Key][]byte{subpartitionedKey: {3 * 42}}, reqctx.uncommittedState.Mutations().Sets)
	require.Equal(t, map[kv.Key]struct{}{}, reqctx.uncommittedState.Mutations().Dels)

	v = s.Get("x")
	require.Equal(t, []byte{3 * 42}, v)
}

func TestIterate(t *testing.T) {
	db := mapdb.NewMapDB()
	cs := state.NewStoreWithUniqueWriteMutex(db)
	origin.InitChain(0, cs, nil, 0)
	latest, err := cs.LatestBlock()
	require.NoError(t, err)
	stateDraft, err := cs.NewStateDraft(time.Now(), latest.L1Commitment())
	require.NoError(t, err)

	hname := isc.Hn("test")

	reqctx := &requestContext{
		vm: &vmContext{
			task:       &vm.VMTask{},
			stateDraft: stateDraft,
		},
		uncommittedState: buffered.NewBufferedKVStore(stateDraft),
		callStack:        []*callContext{{contract: hname}},
	}
	s := reqctx.contractStateWithGasBurn()

	s.Set("xy1", []byte{42})
	s.Set("xy2", []byte{42 * 2})

	arr := make([][]byte, 0)
	s.IterateSorted("xy", func(k kv.Key, v []byte) bool {
		require.True(t, strings.HasPrefix(string(k), "xy"))
		arr = append(arr, v)
		return true
	})
	require.EqualValues(t, 2, len(arr))
	require.Equal(t, []byte{42}, arr[0])
	require.Equal(t, []byte{42 * 2}, arr[1])
}
