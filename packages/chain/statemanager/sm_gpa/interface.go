package sm_gpa

import (
	"time"

	"github.com/nnikolash/wasp-types-exported/packages/chain/statemanager/sm_snapshots"
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/state"
)

type StateManagerOutput interface {
	addBlockCommitted(uint32, *state.L1Commitment)
	TakeBlocksCommitted() []sm_snapshots.SnapshotInfo
	addBlocksToCommit([]*state.L1Commitment)
	TakeNextInputs() []gpa.Input
}

type SnapshotExistsFun func(uint32, *state.L1Commitment) bool

type blockRequestCallback interface {
	isValid() bool
	requestCompleted()
}

type blockFetcher interface {
	getStateIndex() uint32
	getCommitment() *state.L1Commitment
	getCallbacksCount() int
	notifyFetched() []blockFetcher // notifies waiting callbacks of this fetcher and returns all related fetchers
	addCallback(blockRequestCallback)
	addRelatedFetcher(blockFetcher)
	cleanCallbacks()
}

type blockFetchers interface {
	getSize() int
	getCallbacksCount() int
	addFetcher(blockFetcher)
	takeFetcher(*state.L1Commitment) blockFetcher
	addCallback(*state.L1Commitment, blockRequestCallback) bool
	addRelatedFetcher(*state.L1Commitment, blockFetcher) bool
	getCommitments() []*state.L1Commitment
	cleanCallbacks()
}

type blockFetchersMetrics interface {
	inc()
	dec()
	duration(time.Duration)
}
