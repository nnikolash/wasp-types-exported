package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/clients/chainclient"
	"github.com/nnikolash/wasp-types-exported/contracts/native/inccounter"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

func TestMissingRequests(t *testing.T) {
	clu := newCluster(t, waspClusterOpts{nNodes: 4})
	cmt := []int{0, 1, 2, 3}
	threshold := uint16(4)
	addr, err := clu.RunDKG(cmt, threshold)
	require.NoError(t, err)

	chain, err := clu.DeployChain(clu.Config.AllNodes(), cmt, threshold, addr)
	require.NoError(t, err)
	chainID := chain.ChainID

	chEnv := newChainEnv(t, clu, chain)
	chEnv.deployNativeIncCounterSC()

	waitUntil(t, chEnv.contractIsDeployed(), clu.Config.AllNodes(), 30*time.Second)

	userWallet, _, err := chEnv.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	// deposit funds before sending the off-ledger request
	chClient := chainclient.New(clu.L1Client(), clu.WaspClient(0), chainID, userWallet)
	reqTx, err := chClient.DepositFunds(100)
	require.NoError(t, err)
	_, err = chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(chainID, reqTx, false, 30*time.Second)
	require.NoError(t, err)

	// TODO: Validate offleder logic
	// send off-ledger request to all nodes except #3
	req := isc.NewOffLedgerRequest(chainID, nativeIncCounterSCHname, inccounter.FuncIncCounter.Hname(), dict.Dict{}, 0, gas.LimitsDefault.MaxGasPerRequest).Sign(userWallet)

	_, err = clu.WaspClient(0).RequestsApi.OffLedger(context.Background()).OffLedgerRequest(apiclient.OffLedgerRequest{
		ChainId: chainID.String(),
		Request: iotago.EncodeHex(req.Bytes()),
	}).Execute()
	require.NoError(t, err)

	//------
	// send a dummy request to node #3, so that it proposes a batch and the consensus hang is broken
	req2 := isc.NewOffLedgerRequest(chainID, isc.Hn("foo"), isc.Hn("bar"), nil, 1, gas.LimitsDefault.MaxGasPerRequest).Sign(userWallet)

	_, err = clu.WaspClient(0).RequestsApi.OffLedger(context.Background()).OffLedgerRequest(apiclient.OffLedgerRequest{
		ChainId: chainID.String(),
		Request: iotago.EncodeHex(req2.Bytes()),
	}).Execute()
	require.NoError(t, err)
	//-------

	// expect request to be successful, as node #3 must ask for the missing request from other nodes
	waitUntil(t, chEnv.counterEquals(43), clu.Config.AllNodes(), 30*time.Second)
}
