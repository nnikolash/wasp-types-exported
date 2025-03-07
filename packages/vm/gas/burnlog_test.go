package gas_test

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

func TestBurnLogSerialization(t *testing.T) {
	var burnLog gas.BurnLog
	burnLog.Records = []gas.BurnRecord{
		{
			Code:      gas.BurnCodeCallTargetNotFound,
			GasBurned: 10,
		},
		{
			Code:      gas.BurnCodeUtilsHashingSha3,
			GasBurned: 80,
		},
	}
	rwutil.ReadWriteTest(t, &burnLog, new(gas.BurnLog))
}
