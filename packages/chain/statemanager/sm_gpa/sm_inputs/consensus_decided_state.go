package sm_inputs

import (
	"context"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/transaction"
)

type ConsensusDecidedState struct {
	context      context.Context
	stateIndex   uint32
	l1Commitment *state.L1Commitment
	resultCh     chan<- state.State
}

var _ gpa.Input = &ConsensusDecidedState{}

func NewConsensusDecidedState(ctx context.Context, aliasOutput *isc.AliasOutputWithID) (*ConsensusDecidedState, <-chan state.State) {
	commitment, err := transaction.L1CommitmentFromAliasOutput(aliasOutput.GetAliasOutput())
	if err != nil {
		panic("Cannot make L1 commitment from alias output")
	}
	resultChannel := make(chan state.State, 1)
	return &ConsensusDecidedState{
		context:      ctx,
		stateIndex:   aliasOutput.GetStateIndex(),
		l1Commitment: commitment,
		resultCh:     resultChannel,
	}, resultChannel
}

func (cdsT *ConsensusDecidedState) GetStateIndex() uint32 {
	return cdsT.stateIndex
}

func (cdsT *ConsensusDecidedState) GetL1Commitment() *state.L1Commitment {
	return cdsT.l1Commitment
}

func (cdsT *ConsensusDecidedState) IsValid() bool {
	return cdsT.context.Err() == nil
}

func (cdsT *ConsensusDecidedState) Respond(theState state.State) {
	if cdsT.IsValid() && !cdsT.IsResultChClosed() {
		cdsT.resultCh <- theState
		cdsT.closeResultCh()
	}
}

func (cdsT *ConsensusDecidedState) IsResultChClosed() bool {
	return cdsT.resultCh == nil
}

func (cdsT *ConsensusDecidedState) closeResultCh() {
	close(cdsT.resultCh)
	cdsT.resultCh = nil
}
