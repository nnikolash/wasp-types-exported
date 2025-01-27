// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"context"

	"github.com/iotaledger/hive.go/logger"
	"github.com/nnikolash/wasp-types-exported/packages/chain/statemanager"
	"github.com/nnikolash/wasp-types-exported/packages/chain/statemanager/sm_gpa/sm_inputs"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/state"
)

// Tracks a single chain of state transitions. We will have 2 instances of it:
//   - one for tracking the active state. It is needed for mempool to clear the requests.
//   - one for the committed state to await for committed request receipts.
type StateTracker interface {
	//
	// The main functions provided by this component.
	TrackAliasOutput(ao *isc.AliasOutputWithID, strict bool)
	AwaitRequestReceipt(query *awaitReceiptReq)
	//
	// The following 2 functions are only to move the channel receive loop to the main ChainNode thread.
	ChainNodeAwaitStateMgrCh() <-chan *sm_inputs.ChainFetchStateDiffResults
	ChainNodeStateMgrResponse(*sm_inputs.ChainFetchStateDiffResults)
}

type StateTrackerStepCB = func(st state.State, from, till *isc.AliasOutputWithID, added, removed []state.Block)

type stateTrackerImpl struct {
	ctx                    context.Context
	stateMgr               statemanager.StateMgr
	haveLatestCB           StateTrackerStepCB
	haveAOState            state.State
	haveAO                 *isc.AliasOutputWithID // We have a state ready for this AO.
	nextAO                 *isc.AliasOutputWithID // For this state a query was made, but the response not received yet.
	nextAOCancel           context.CancelFunc     // Cancel for a context used to query for the nextAO state.
	nextAOWaitCh           <-chan *sm_inputs.ChainFetchStateDiffResults
	awaitReceipt           AwaitReceipt
	metricWantStateIndexCB func(uint32)
	metricHaveStateIndexCB func(uint32)
	log                    *logger.Logger
}

var _ StateTracker = &stateTrackerImpl{}

func NewStateTracker(
	ctx context.Context,
	stateMgr statemanager.StateMgr,
	haveLatestCB StateTrackerStepCB,
	metricWantStateIndexCB func(uint32),
	metricHaveStateIndexCB func(uint32),
	log *logger.Logger,
) StateTracker {
	return &stateTrackerImpl{
		ctx:                    ctx,
		stateMgr:               stateMgr,
		haveLatestCB:           haveLatestCB,
		haveAOState:            nil,
		haveAO:                 nil,
		nextAO:                 nil,
		nextAOCancel:           nil,
		nextAOWaitCh:           nil,
		awaitReceipt:           NewAwaitReceipt(AwaitReceiptCleanupEvery, log),
		metricWantStateIndexCB: metricWantStateIndexCB,
		metricHaveStateIndexCB: metricHaveStateIndexCB,
		log:                    log,
	}
}

func (sti *stateTrackerImpl) TrackAliasOutput(ao *isc.AliasOutputWithID, strict bool) {
	sti.log.Debugf("TrackAliasOutput[strict=%v], ao=%v, haveAO=%v, nextAO=%v", strict, ao, sti.haveAO, sti.nextAO)
	if !strict && sti.haveAO != nil && sti.haveAO.GetStateIndex() >= ao.GetStateIndex() {
		return
	}
	if ao.Equals(sti.nextAO) {
		return
	}
	sti.metricWantStateIndexCB(ao.GetStateIndex())
	if ao.Equals(sti.haveAO) {
		sti.nextAO = sti.haveAO // All done, state is already received.
		sti.cancelQuery()       // Cancel the request, if pending.
		return
	}
	nextAOCtx, nextAOCancel := context.WithCancel(sti.ctx)
	sti.nextAO = ao
	sti.nextAOCancel = nextAOCancel
	sti.nextAOWaitCh = sti.stateMgr.ChainFetchStateDiff(nextAOCtx, sti.haveAO, sti.nextAO)
}

func (sti *stateTrackerImpl) AwaitRequestReceipt(query *awaitReceiptReq) {
	sti.log.Debugf("AwaitRequestReceipt, query.requestID=%v", query.requestID)
	sti.awaitReceipt.Await(query)
}

// To be used in the select loop at the chain node.
func (sti *stateTrackerImpl) ChainNodeAwaitStateMgrCh() <-chan *sm_inputs.ChainFetchStateDiffResults {
	return sti.nextAOWaitCh
}

// This is assumed to be called right after the `ChainNodeAwaitStateMgrCh()`,
// thus no additional checks are present here.
func (sti *stateTrackerImpl) ChainNodeStateMgrResponse(results *sm_inputs.ChainFetchStateDiffResults) {
	sti.cancelQuery()
	newState := results.GetNewState()
	sti.log.Debugf(
		"Have latest state for %v, state.BlockIndex=%v, state.trieRoot=%v, previous=%v, |blocksAdded|=%v, |blockRemoved|=%v",
		sti.nextAO, newState.BlockIndex(), newState.TrieRoot(), sti.haveAO, len(results.GetAdded()), len(results.GetRemoved()),
	)
	sti.haveLatestCB(newState, sti.haveAO, sti.nextAO, results.GetAdded(), results.GetRemoved())
	sti.haveAO = sti.nextAO
	sti.haveAOState = newState
	sti.metricHaveStateIndexCB(newState.BlockIndex())
	sti.awaitReceipt.ConsiderState(newState, results.GetAdded())
}

func (sti *stateTrackerImpl) cancelQuery() {
	if sti.nextAOCancel == nil {
		return
	}
	sti.nextAOCancel()
	sti.nextAOCancel = nil
	sti.nextAOWaitCh = nil
}
