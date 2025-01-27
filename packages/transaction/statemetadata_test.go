package transaction_test

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/transaction"
	"github.com/nnikolash/wasp-types-exported/packages/util"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

func TestStateMetadataSerialization(t *testing.T) {
	s := transaction.NewStateMetadata(
		state.PseudoRandL1Commitment(),
		&gas.FeePolicy{
			GasPerToken: util.Ratio32{
				A: 1,
				B: 2,
			},
			EVMGasRatio: util.Ratio32{
				A: 3,
				B: 4,
			},
			ValidatorFeeShare: 5,
		},
		6,
		"https://iota.org",
	)
	rwutil.BytesTest(t, s, transaction.StateMetadataFromBytes)
}
