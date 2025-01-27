package test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/solo"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/utxodb"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/common"
)

func TestOffLedger(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{
		AutoAdjustStorageDeposit: true,
		GasBurnLogEnabled:        true,
	})
	chain := env.NewChain()

	// create a wallet with some base tokens on L1:
	userWallet, userAddress := env.NewKeyPairWithFunds(env.NewSeedFromIndex(0))
	env.AssertL1BaseTokens(userAddress, utxodb.FundsFromFaucetAmount)
	chain.DepositBaseTokensToL2(env.L1BaseTokens(userAddress), userWallet)

	req := isc.NewOffLedgerRequest(chain.ID(), accounts.Contract.Hname(), accounts.FuncDeposit.Hname(), dict.New(), 0, math.MaxUint64)
	altReq := isc.NewImpersonatedOffLedgerRequest(req.(*isc.OffLedgerRequestData)).
		WithSenderAddress(userWallet.Address())

	rec, err := common.EstimateGas(chain, altReq)

	require.NoError(t, err)
	require.NotNil(t, rec)
	require.Greater(t, rec.GasFeeCharged, uint64(0))
}
