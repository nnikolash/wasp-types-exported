package sm_snapshots

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/nnikolash/wasp-types-exported/packages/chain/statemanager/sm_gpa/sm_gpa_utils"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testlogger"
)

func TestWriteReadDifferentStores(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()

	var err error
	numberOfBlocks := 10
	factory := sm_gpa_utils.NewBlockFactory(t)
	blocks := factory.GetBlocks(numberOfBlocks, 1)
	lastBlock := blocks[numberOfBlocks-1]
	lastCommitment := lastBlock.L1Commitment()
	snapshotInfo := NewSnapshotInfo(blocks[numberOfBlocks-1].StateIndex(), lastCommitment)
	snapshotterOrig := newSnapshotter(factory.GetStore())
	fileName := "TestWriteReadDifferentStores.snap"
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	require.NoError(t, err)
	err = snapshotterOrig.storeSnapshot(snapshotInfo, f)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)

	store := state.NewStoreWithUniqueWriteMutex(mapdb.NewMapDB())
	snapshotterNew := newSnapshotter(store)
	f, err = os.Open(fileName)
	require.NoError(t, err)
	err = snapshotterNew.loadSnapshot(snapshotInfo, f)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	err = os.Remove(fileName)
	require.NoError(t, err)

	sm_gpa_utils.CheckBlockInStore(t, store, lastBlock)
	sm_gpa_utils.CheckStateInStores(t, factory.GetStore(), store, lastCommitment)
}
