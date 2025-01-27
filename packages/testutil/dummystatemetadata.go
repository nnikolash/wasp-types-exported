package testutil

import (
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/transaction"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

func DummyStateMetadata(commitment *state.L1Commitment) *transaction.StateMetadata {
	return transaction.NewStateMetadata(
		commitment,
		gas.DefaultFeePolicy(),
		0,
		"",
	)
}
