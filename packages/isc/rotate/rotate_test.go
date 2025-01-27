package rotate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testkey"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

func TestBasicRotateRequest(t *testing.T) {
	kp, addr := testkey.GenKeyAddr()
	req := NewRotateRequestOffLedger(isc.RandomChainID(), addr, kp, gas.LimitsDefault.MaxGasPerRequest)
	require.True(t, IsRotateStateControllerRequest(req))
}
