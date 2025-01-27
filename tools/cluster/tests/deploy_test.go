package tests

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/clients/apiextensions"
	"github.com/nnikolash/wasp-types-exported/clients/chainclient"
	"github.com/nnikolash/wasp-types-exported/contracts/native/inccounter"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/corecontracts"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/root"
)

// executed in cluster_test.go
func testDeployChain(t *testing.T, env *ChainEnv) {
	chainID, chainOwnerID := env.getChainInfo()
	require.EqualValues(t, chainID, env.Chain.ChainID)
	require.EqualValues(t, chainOwnerID, isc.NewAgentID(env.Chain.OriginatorAddress()))
	t.Logf("--- chainID: %s", chainID.String())
	t.Logf("--- chainOwnerID: %s", chainOwnerID.String())

	env.checkCoreContracts()
	env.checkRootsOutside()
	for _, i := range env.Chain.CommitteeNodes {
		blockIndex, err := env.Chain.BlockIndex(i)
		require.NoError(t, err)
		require.Greater(t, blockIndex, uint32(1))

		contractRegistry, err := env.Chain.ContractRegistry(i)
		require.NoError(t, err)
		require.EqualValues(t, len(corecontracts.All), len(contractRegistry))
	}
}

// executed in cluster_test.go
func testDeployContractOnly(t *testing.T, env *ChainEnv) {
	env.deployNativeIncCounterSC()

	// test calling root.FuncFindContractByName view function using client
	ret, err := apiextensions.CallView(
		context.Background(),
		env.Chain.Cluster.WaspClient(),
		env.Chain.ChainID.String(),
		apiclient.ContractCallViewRequest{
			ContractHName: root.Contract.Hname().String(),
			FunctionHName: root.ViewFindContract.Hname().String(),
			Arguments: apiextensions.DictToAPIJsonDict(dict.Dict{
				root.ParamHname: isc.Hn(nativeIncCounterSCName).Bytes(),
			}),
		})

	require.NoError(t, err)
	recb := ret.Get(root.ParamContractRecData)
	_, err = root.ContractRecordFromBytes(recb)
	require.NoError(t, err)
}

// executed in cluster_test.go
func testDeployContractAndSpawn(t *testing.T, env *ChainEnv) {
	env.deployNativeIncCounterSC()

	hname := isc.Hn(nativeIncCounterSCName)

	nameNew := "spawnedContract"
	hnameNew := isc.Hn(nameNew)
	// send 'spawn' request to the SC which was just deployed
	par := chainclient.NewPostRequestParams(
		inccounter.VarName, nameNew,
	).WithBaseTokens(100)
	tx, err := env.Chain.OriginatorClient().Post1Request(hname, inccounter.FuncSpawn.Hname(), *par)
	require.NoError(t, err)

	receipts, err := env.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessed(env.Chain.ChainID, tx, false, 30*time.Second)
	require.NoError(t, err)
	require.Len(t, receipts, 1)

	env.checkCoreContracts()
	for _, i := range env.Chain.CommitteeNodes {
		blockIndex, err := env.Chain.BlockIndex(i)
		require.NoError(t, err)
		require.Greater(t, blockIndex, uint32(2))

		contractRegistry, err := env.Chain.ContractRegistry(i)
		require.NoError(t, err)
		require.EqualValues(t, len(corecontracts.All)+2, len(contractRegistry))

		cr, ok := lo.Find(contractRegistry, func(item apiclient.ContractInfoResponse) bool {
			return item.HName == hnameNew.String()
		})
		require.True(t, ok)
		require.NotNil(t, cr)

		require.EqualValues(t, nameNew, cr.Name)

		counterValue, err := env.Chain.GetCounterValue(hname, i)
		require.NoError(t, err)
		require.EqualValues(t, 42, counterValue)
	}
}
