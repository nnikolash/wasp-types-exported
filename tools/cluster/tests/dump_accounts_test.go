package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/solo"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/utxodb"
)

func testDumpAccounts(t *testing.T, env *ChainEnv) {
	// create 10 accounts with funds on-chain
	accs := make([]string, 0, 100)
	for i := 0; i < 5; i++ {
		// 5 L1 accounts
		keyPair, addr, err := env.Clu.NewKeyPairWithFunds()
		require.NoError(t, err)
		env.DepositFunds(10*isc.Million, keyPair)
		accs = append(accs, addr.Bech32(parameters.L1().Protocol.Bech32HRP))
	}

	for i := 0; i < 5; i++ {
		// 5 EVM accounts
		_, evmAddr := solo.NewEthereumAccount()
		keyPair, _, err := env.Clu.NewKeyPairWithFunds()
		require.NoError(t, err)
		evmAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, evmAddr)
		env.TransferFundsTo(isc.NewAssetsBaseTokens(utxodb.FundsFromFaucetAmount-1*isc.Million), nil, keyPair, evmAgentID)
		accs = append(accs, evmAgentID.String())
	}

	resp, err := env.NewChainClient().WaspClient.ChainsApi.DumpAccounts(
		context.Background(),
		env.Chain.ChainID.String(),
	).Execute()
	require.NoError(t, err)
	require.Equal(t, 202, resp.StatusCode)
	time.Sleep(1 * time.Second) // wait for the file to be produced

	path := filepath.Join(env.Clu.NodeDataPath(0), "waspdb", "account_dumps", env.Chain.ChainID.String())
	entries, err := os.ReadDir(path)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	contents, err := os.ReadFile(filepath.Join(path, entries[0].Name()))
	require.NoError(t, err)
	// assert all accounts are present in the dump
	for _, acc := range accs {
		require.Contains(t, string(contents), acc)
	}
}
