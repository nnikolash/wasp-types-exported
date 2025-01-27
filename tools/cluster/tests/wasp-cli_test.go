package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"

	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testkey"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blocklog"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
	"github.com/nnikolash/wasp-types-exported/packages/vm/vmtypes"
	"github.com/nnikolash/wasp-types-exported/tools/cluster/templates"
)

const file = "inccounter_bg.wasm"

const srcFile = "wasm/" + file

func TestWaspCLINoChains(t *testing.T) {
	w := newWaspCLITest(t)

	out := w.MustRun("address")

	ownerAddr := regexp.MustCompile(`(?m)Address:\s+([[:alnum:]]+)$`).FindStringSubmatch(out[1])[1]
	require.NotEmpty(t, ownerAddr)

	out = w.MustRun("chain", "list", "--node=0", "--node=0")
	require.Contains(t, out[0], "Total 0 chain(s)")
}

func TestWaspAuth(t *testing.T) {
	w := newWaspCLITest(t, waspClusterOpts{
		modifyConfig: func(nodeIndex int, configParams templates.WaspConfigParams) templates.WaspConfigParams {
			configParams.AuthScheme = "jwt"
			return configParams
		},
	})
	_, err := w.Run("chain", "list", "--node=0", "--node=0")
	require.Error(t, err)
	out := w.MustRun("auth", "login", "--node=0", "-u=wasp", "-p=wasp")
	require.Equal(t, "Successfully authenticated", out[1])
	out = w.MustRun("chain", "list", "--node=0", "--node=0")
	require.Contains(t, out[0], "Total 0 chain(s)")
}

func TestZeroGasFee(t *testing.T) {
	w := newWaspCLITest(t)

	const chainName = "chain1"
	committee, quorum := w.ArgCommitteeConfig(0)

	// test chain deploy command
	w.MustRun("chain", "deploy", "--chain="+chainName, committee, quorum, "--block-keep-amount=123", "--node=0")
	w.ActivateChainOnAllNodes(chainName, 0)
	outs, err := w.Run("chain", "info", "--node=0", "--node=0")
	require.NoError(t, err)
	require.Contains(t, outs, "Gas fee: gas units * (100/1)")
	_, err = w.Run("chain", "disable-feepolicy", "--node=0")
	require.NoError(t, err)
	outs, err = w.Run("chain", "info", "--node=0", "--node=0")
	require.NoError(t, err)
	require.Contains(t, outs, "Gas fee: gas units * (0/0)")

	t.Run("send arbitrary EVM tx without funds", func(t *testing.T) {
		ethPvtKey, _ := newEthereumAccount()
		sendDummyEVMTx(t, w, ethPvtKey)
	})

	t.Run("deposit directly to EVM", func(t *testing.T) {
		alternativeAddress := getAddress(w.MustRun("address", "--address-index=1"))
		w.MustRun("send-funds", "-s", alternativeAddress, "base:1000000")
		checkBalance(t, w.MustRun("balance", "--address-index=1"), 1000000)
		_, eth := newEthereumAccount()
		w.MustRun("chain", "deposit", eth.String(), "base:1000000", "--node=0", "--address-index=1")
		checkBalance(t, w.MustRun("chain", "balance", eth.String(), "--node=0"), 1000000)
	})
}

func checkBalance(t *testing.T, out []string, expected int) {
	t.Helper()
	// regex example: base tokens 1000000
	//				  -----  ------token  amount-----  ------base   1364700
	r := regexp.MustCompile(`.*(?i:base)\s*(?i:tokens)?:*\s*(\d+).*`).FindStringSubmatch(strings.Join(out, ""))
	if r == nil {
		panic("couldn't check balance")
	}
	amount, err := strconv.Atoi(r[1])
	require.NoError(t, err)
	require.EqualValues(t, expected, amount)
}

func getAddress(out []string) string {
	r := regexp.MustCompile(`.*Address:\s+(\w*).*`).FindStringSubmatch(strings.Join(out, ""))
	if r == nil {
		panic("couldn't get address")
	}
	return r[1]
}

func TestWaspCLISendFunds(t *testing.T) {
	w := newWaspCLITest(t)

	alternativeAddress := getAddress(w.MustRun("address", "--address-index=1"))

	w.MustRun("send-funds", "-s", alternativeAddress, "base:1000000")
	checkBalance(t, w.MustRun("balance", "--address-index=1"), 1000000)
}

func TestWaspCLIDeposit(t *testing.T) {
	w := newWaspCLITest(t)

	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)

	// fund an alternative address to deposit from (so we can test the fees, since --address-index=0 is the chain owner / default payoutAddress)
	alternativeAddress := getAddress(w.MustRun("address", "--address-index=1"))
	w.MustRun("send-funds", "-s", alternativeAddress, "base:10000000")

	minFee := gas.DefaultFeePolicy().MinFee(nil, parameters.L1().BaseToken.Decimals)
	t.Run("deposit directly to EVM", func(t *testing.T) {
		_, eth := newEthereumAccount()
		w.MustRun("chain", "deposit", eth.String(), "base:1000000", "--node=0", "--address-index=1")
		checkBalance(t, w.MustRun("chain", "balance", eth.String(), "--node=0"), 1000000-int(minFee))
	})

	t.Run("deposit to own account, then to EVM", func(t *testing.T) {
		w.MustRun("chain", "deposit", "base:1000000", "--node=0", "--address-index=1")
		checkBalance(t, w.MustRun("chain", "balance", "--node=0", "--address-index=1"), 1000000-int(minFee))
		_, eth := newEthereumAccount()
		w.MustRun("chain", "deposit", eth.String(), "base:1000000", "--node=0", "--address-index=1")
		checkBalance(t, w.MustRun("chain", "balance", eth.String(), "--node=0", "--address-index=1"), 1000000) // fee will be taken from the sender on-chain balance
		checkBalance(t, w.MustRun("chain", "balance", "--node=0", "--address-index=1"), 1000000-2*int(minFee))
	})

	t.Run("mint and deposit native tokens to an ethereum account", func(t *testing.T) {
		_, eth := newEthereumAccount()
		// create foundry
		tokenScheme := codec.EncodeTokenScheme(&iotago.SimpleTokenScheme{
			MintedTokens:  big.NewInt(0),
			MeltedTokens:  big.NewInt(0),
			MaximumSupply: big.NewInt(1000),
		})

		out := w.PostRequestGetReceipt(
			"accounts", accounts.FuncNativeTokenCreate.Name,
			"string", accounts.ParamTokenScheme, "bytes", iotago.EncodeHex(tokenScheme),
			"-l", "base:1000000",
			"-t", "base:100000000",
			"string", accounts.ParamTokenName, "string", "TEST",
			"string", accounts.ParamTokenTickerSymbol, "string", "TS",
			"string", accounts.ParamTokenDecimals, "uint8", "8",
			"--node=0",
		)
		require.Regexp(t, `.*Error: \(empty\).*`, strings.Join(out, ""))

		// mint 2 native tokens
		foundrySN := "1"
		out = w.PostRequestGetReceipt(
			"accounts", accounts.FuncNativeTokenModifySupply.Name,
			"string", accounts.ParamFoundrySN, "uint32", foundrySN,
			"string", accounts.ParamSupplyDeltaAbs, "bigint", "2",
			"string", accounts.ParamDestroyTokens, "bool", "false",
			"-l", "base:1000000",
			"--off-ledger",
			"--node=0",
		)
		require.Regexp(t, `.*Error: \(empty\).*`, strings.Join(out, ""))

		out = w.MustRun("chain", "balance", "--node=0")
		tokenID := ""
		for _, line := range out {
			if strings.Contains(line, "0x") {
				tokenID = strings.Split(line, " ")[0]
			}
		}

		// withdraw this token to the wasp-cli L1 address
		out = w.PostRequestGetReceipt(
			"accounts", accounts.FuncWithdraw.Name,
			"-l", fmt.Sprintf("base:1000000, %s:2", tokenID),
			"--off-ledger",
			"--node=0",
		)
		require.Regexp(t, `.*Error: \(empty\).*`, strings.Join(out, ""))

		out = w.MustRun("balance")
		for _, line := range out {
			if strings.Contains(line, "0x") {
				balance := strings.TrimSpace(strings.Split(line, ":")[1])
				require.Equal(t, balance, "2")
			}
		}

		// deposit the native token to the chain (to an ethereum account)
		w.MustRun(
			"chain", "deposit", eth.String(),
			fmt.Sprintf("%s:1", tokenID),
			"--adjust-storage-deposit",
			"--node=0",
		)
		out = w.MustRun("chain", "balance", eth.String(), "--node=0")
		require.Contains(t, strings.Join(out, ""), tokenID)

		// deposit the native token to the chain (to the cli account)
		w.MustRun(
			"chain", "deposit",
			fmt.Sprintf("%s:1", tokenID),
			"--adjust-storage-deposit",
			"--node=0",
		)
		out = w.MustRun("chain", "balance", "--node=0")
		require.Contains(t, strings.Join(out, ""), tokenID)
		// no token balance on L1
		out = w.MustRun("balance")
		require.NotContains(t, strings.Join(out, ""), tokenID)
	})
}

func TestWaspCLIUnprocessableRequest(t *testing.T) {
	w := newWaspCLITest(t)

	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)
	regex4NTsBalance := regexp.MustCompile(`.*(0x.{76}).*(0x.{76}).*(0x.{76}).*(0x.{76}).*`)

	createFoundries := func(nFoundries int) []string {
		// create N foundries and mints 1 tokens from each foundry
		tokenScheme := codec.EncodeTokenScheme(&iotago.SimpleTokenScheme{
			MintedTokens:  big.NewInt(0),
			MeltedTokens:  big.NewInt(0),
			MaximumSupply: big.NewInt(1),
		})
		for i := 1; i <= nFoundries; i++ {
			// create foundry
			out := w.PostRequestGetReceipt(
				"accounts", accounts.FuncNativeTokenCreate.Name,
				"string", accounts.ParamTokenScheme, "bytes", iotago.EncodeHex(tokenScheme),
				"-l", "base:1000000",
				"-t", "base:100000000",
				"string", accounts.ParamTokenName, "string", "TEST",
				"string", accounts.ParamTokenTickerSymbol, "string", "TS",
				"string", accounts.ParamTokenDecimals, "uint8", "8",
				"--node=0",
			)
			require.Regexp(t, `.*Error: \(empty\).*`, strings.Join(out, ""))

			// mint 1 native token
			out = w.PostRequestGetReceipt(
				"accounts", accounts.FuncNativeTokenModifySupply.Name,
				"string", accounts.ParamFoundrySN, "uint32", fmt.Sprintf("%d", i),
				"string", accounts.ParamSupplyDeltaAbs, "bigint", "1",
				"string", accounts.ParamDestroyTokens, "bool", "false",
				"-l", "base:1000000",
				"--off-ledger",
				"--node=0",
			)
			require.Regexp(t, `.*Error: \(empty\).*`, strings.Join(out, ""))
		}

		out := w.MustRun("chain", "balance", "--node=0")
		ntIDs := regex4NTsBalance.FindStringSubmatch(strings.Join(out, ""))
		ntIDs = ntIDs[1:]

		// withdraw all tokens to the target L1 address
		out = w.PostRequestGetReceipt(
			"accounts", accounts.FuncWithdraw.Name,
			"-l", fmt.Sprintf("base:1000000, %s:1, %s:1, %s:1, %s:1", ntIDs[0], ntIDs[1], ntIDs[2], ntIDs[3]),
			"--off-ledger",
			"--node=0",
		)
		require.Regexp(t, `.*Error: \(empty\).*`, strings.Join(out, ""))

		return ntIDs
	}

	// create 4 foundries, mint 1 token of each
	ntIDs := createFoundries(4)

	// send those tokens to another address (that doesn't have on-chain balance)
	alternativeAddress := getAddress(w.MustRun("address", "--address-index=1"))
	w.MustRun("send-funds", "-s", alternativeAddress, fmt.Sprintf("base:10000000, %s:1, %s:1, %s:1, %s:1", ntIDs[0], ntIDs[1], ntIDs[2], ntIDs[3]))
	out := w.MustRun("balance", "--address-index=1")
	require.True(t, regex4NTsBalance.Match([]byte(strings.Join(out, ""))))

	go func() {
		// wait some time and post another request so that a block is produced (a block won't be produced with only unprocessable requests)
		time.Sleep(1 * time.Second)
		w.MustRun("chain", "deposit", "base:1000000", "--node=0")
	}()
	// try to deposit them without enough funds (only send minSD) // THIS IS THIS UNPROCESSABLE REQUEST
	out = w.MustRun(
		"chain", "deposit",
		fmt.Sprintf("%s:1, %s:1, %s:1, %s:1", ntIDs[0], ntIDs[1], ntIDs[2], ntIDs[3]),
		"--adjust-storage-deposit",
		"--node=0",
		"--address-index=1",
	)
	reqID := findRequestIDInOutput(out)
	require.NotEmpty(t, reqID)

	// native tokens don't appear on-chain, nor on L1
	out = w.MustRun("balance", "--address-index=1")
	require.False(t, regex4NTsBalance.Match([]byte(strings.Join(out, ""))))

	out = w.MustRun("chain", "balance", "--node=0", "--address-index=1")
	require.False(t, regex4NTsBalance.Match([]byte(strings.Join(out, ""))))

	// check "request to retry" exists
	out = w.MustRun("chain", "call-view",
		blocklog.Contract.Name, blocklog.ViewHasUnprocessable.Name,
		"string", blocklog.ParamRequestID, "bytes", reqID,
		"--node=0")
	require.Regexp(t, `"value":"0x01"`, strings.Join(out, ""))

	// send a "retry" request that deposits enough funds to process the initial NTs deposit
	out = w.PostRequestGetReceipt(
		blocklog.Contract.Name, blocklog.FuncRetryUnprocessable.Name,
		"string", blocklog.ParamRequestID, "bytes", reqID,
		"-t", "base:1000000",
		"--node=0",
		"--address-index=1",
	)
	require.Regexp(t, `.*Error: \(empty\).*`, strings.Join(out, ""))

	// check "request to retry" has left the "unprocessable list"
	out = w.MustRun("chain", "call-view",
		blocklog.Contract.Name, blocklog.ViewHasUnprocessable.Name,
		"string", blocklog.ParamRequestID, "bytes", reqID,
		"--node=0")
	require.Regexp(t, `"value":"0x00"`, strings.Join(out, ""))

	// native tokens now appear on-chain, nor on L1
	out = w.MustRun("chain", "balance", "--node=0", "--address-index=1")
	require.Len(t, regex4NTsBalance.FindStringSubmatch(strings.Join(out, "")), 5)
}

func TestWaspCLIContract(t *testing.T) {
	w := newWaspCLITest(t)

	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)

	// for running off-ledger requests
	w.MustRun("chain", "deposit", "base:10000000", "--node=0")

	vmtype := vmtypes.WasmTime
	name := "inccounter"
	description := "inccounter SC"
	w.CopyFile(srcFile)

	// test chain deploy-contract command
	w.MustRun("chain", "deploy-contract", vmtype, name, description, file,
		"string", "counter", "int64", "42",
		"--node=0",
	)

	out := w.MustRun("chain", "list-contracts", "--node=0")
	found := false
	for _, s := range out {
		if strings.Contains(s, name) {
			found = true
			break
		}
	}
	require.True(t, found)

	checkCounter := func(n int) {
		// test chain call-view command
		out = w.MustRun("chain", "call-view", name, "getCounter", "--node=0")
		out = w.MustPipe(out, "decode", "string", "counter", "int")
		require.Regexp(t, fmt.Sprintf(`(?m)counter:\s+%d$`, n), out[0])
	}

	checkCounter(42)

	// test chain post-request command
	w.MustRun("chain", "post-request", "-s", name, "increment", "--node=0")
	checkCounter(43)

	// include a funds transfer
	w.MustRun("chain", "post-request", "-s", name, "increment", "--transfer=base:10000000", "--node=0")
	checkCounter(44)

	// test off-ledger request
	w.MustRun("chain", "post-request", "-s", name, "increment", "--off-ledger", "--node=0")
	checkCounter(45)

	// include an allowance transfer
	w.MustRun("chain", "post-request", "-s", name, "increment", "--transfer=base:10000000", "--allowance=base:10000000", "--node=0")
	checkCounter(46)
}

func findRequestIDInOutput(out []string) string {
	for _, line := range out {
		m := regexp.MustCompile(`(?m)\(check result with: wasp-cli chain request ([-\w]+)\)$`).FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}
		return m[1]
	}
	return ""
}

func TestWaspCLIBlockLog(t *testing.T) {
	w := newWaspCLITest(t)

	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)

	out := w.MustRun("chain", "deposit", "base:100", "--node=0")
	reqID := findRequestIDInOutput(out)
	require.NotEmpty(t, reqID)

	out = w.MustRun("chain", "block", "--node=0")
	require.Equal(t, "Block index: 1", out[0])
	found := false
	for _, line := range out {
		if strings.Contains(line, reqID) {
			found = true
			break
		}
	}
	require.True(t, found)

	out = w.MustRun("chain", "block", "1", "--node=0")
	require.Equal(t, "Block index: 1", out[0])

	out = w.MustRun("chain", "request", reqID, "--node=0")
	t.Log(out)
	found = false
	for _, line := range out {
		if strings.Contains(line, "Error: (empty)") {
			found = true
			break
		}
	}
	require.True(t, found)

	// try an unsuccessful request (missing params)
	out = w.MustRun("chain", "post-request", "-s", "root", "deployContract", "string", "foo", "string", "bar", "--node=0")
	reqID = findRequestIDInOutput(out)
	require.NotEmpty(t, reqID)

	out = w.MustRun("chain", "request", reqID, "--node=0")

	found = false
	for _, line := range out {
		if strings.Contains(line, "Error: ") {
			found = true
			require.Regexp(t, `cannot decode`, line)
			break
		}
	}
	require.True(t, found)

	found = false
	for _, line := range out {
		if strings.Contains(line, "foo") {
			found = true
			require.Contains(t, line, iotago.EncodeHex([]byte("bar")))
			break
		}
	}
	require.True(t, found)
}

func TestWaspCLIRejoinChain(t *testing.T) {
	w := newWaspCLITest(t)

	// make sure deploying with a bad quorum breaks
	require.Panics(
		t,
		func() {
			w.MustRun("chain", "deploy", "--chain=chain1", "--peers=0,1,2,3,4,5", "--quorum=4", "--node=0")
			w.ActivateChainOnAllNodes("chain1", 0)
		})

	chainName := "chain1"

	committee, quorum := w.ArgCommitteeConfig(0)

	// test chain deploy command
	w.MustRun("chain", "deploy", "--chain="+chainName, committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes(chainName, 0)

	var chainID string
	for _, idx := range w.Cluster.AllNodes() {
		// test chain info command
		chainID = w.ChainID(idx)
		require.NotEmpty(t, chainID)
		t.Logf("Chain ID: %s", chainID)
	}

	// test chain list command
	for _, idx := range w.Cluster.AllNodes() {
		out := w.MustRun("chain", "list", fmt.Sprintf("--node=%d", idx))
		require.Contains(t, out[0], "Total 1 chain(s)")
		require.Contains(t, out[4], chainID)
	}

	for _, idx := range w.Cluster.AllNodes() {
		// deactivate chain and check that the chain was deactivated
		w.MustRun("chain", "deactivate", fmt.Sprintf("--node=%d", idx))
		out := w.MustRun("chain", "list", fmt.Sprintf("--node=%d", idx))
		require.Contains(t, out[0], "Total 1 chain(s)")
		require.Contains(t, out[4], chainID)

		chOut := strings.Fields(out[4])
		active, _ := strconv.ParseBool(chOut[1])
		require.False(t, active)
	}

	for _, idx := range w.Cluster.AllNodes() {
		// activate chain and check that it was activated
		w.MustRun("chain", "activate", fmt.Sprintf("--node=%d", idx))
		out := w.MustRun("chain", "list", fmt.Sprintf("--node=%d", idx))
		require.Contains(t, out[0], "Total 1 chain(s)")
		require.Contains(t, out[4], chainID)

		chOut := strings.Fields(out[4])
		active, _ := strconv.ParseBool(chOut[1])
		require.True(t, active)
	}
}

func TestWaspCLILongParam(t *testing.T) {
	w := newWaspCLITest(t)

	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)
	w.MustRun("chain", "deposit", "base:1000000", "--node=0")

	veryLongTokenName := strings.Repeat("A", 10_000)

	errMsg := "slice length is too long"
	defer func() {
		if r := recover(); r != nil {
			errStr := fmt.Sprintf("%s", r)
			if !strings.Contains(errStr, errMsg) {
				t.FailNow()
			}
		}
	}()

	w.CreateL2NativeToken(&iotago.SimpleTokenScheme{
		MaximumSupply: big.NewInt(1000000),
		MeltedTokens:  big.NewInt(0),
		MintedTokens:  big.NewInt(0),
	}, veryLongTokenName, "TST", 8)

	// The code should not reach here. CreateL2NativeToken should panic as the args are too long.
	// This is caught by the deferred recover.
	t.FailNow()
}

func TestWaspCLITrustListImport(t *testing.T) {
	w := newWaspCLITest(t, waspClusterOpts{
		nNodes:  4,
		dirName: "wasp-cluster-initial",
	})

	w2 := newWaspCLITest(t, waspClusterOpts{
		nNodes:  2,
		dirName: "wasp-cluster-new-gov",
		modifyConfig: func(nodeIndex int, configParams templates.WaspConfigParams) templates.WaspConfigParams {
			// avoid port conflicts when running everything on localhost
			configParams.APIPort += 100
			configParams.MetricsPort += 100
			configParams.PeeringPort += 100
			configParams.ProfilingPort += 100
			return configParams
		},
	})

	// set cluster2/node0 to trust all nodes from cluster 1
	for _, nodeIndex := range w.Cluster.Config.AllNodes() {
		peeringInfoOutput := w.MustRun("peering", "info", fmt.Sprintf("--node=%d", nodeIndex))
		pubKey := regexp.MustCompile(`PubKey:\s+([[:alnum:]]+)$`).FindStringSubmatch(peeringInfoOutput[0])[1]
		peeringURL := regexp.MustCompile(`PeeringURL:\s+(.+)$`).FindStringSubmatch(peeringInfoOutput[1])[1]
		w2.MustRun("peering", "trust", fmt.Sprintf("x%d", nodeIndex), pubKey, peeringURL, "--node=0")
	}

	// import the trust from cluster2/node0 to cluster2/node1
	trustedFile0, err := os.CreateTemp("", "tmp-trusted-peers.*.json")
	require.NoError(t, err)
	defer os.Remove(trustedFile0.Name())
	w2.MustRun("peering", "export-trusted", "--node=0", "--peers=x0,x1,x2,x3", "-o="+trustedFile0.Name())
	w2.MustRun("peering", "import-trusted", trustedFile0.Name(), "--node=1")

	// export the trusted nodes from cluster2/node1 and assert the expected result
	trustedFile1, err := os.CreateTemp("", "tmp-trusted-peers.*.json")
	require.NoError(t, err)
	defer os.Remove(trustedFile1.Name())
	w2.MustRun("peering", "export-trusted", "--peers=x0,x1,x2,x3", "--node=1", "-o="+trustedFile1.Name())

	trustedBytes0, err := io.ReadAll(trustedFile0)
	require.NoError(t, err)
	trustedBytes1, err := io.ReadAll(trustedFile1)
	require.NoError(t, err)

	var trustedList0 []apiclient.PeeringNodeIdentityResponse
	require.NoError(t, json.Unmarshal(trustedBytes0, &trustedList0))

	var trustedList1 []apiclient.PeeringNodeIdentityResponse
	require.NoError(t, json.Unmarshal(trustedBytes1, &trustedList1))

	require.Equal(t, len(trustedList0), len(trustedList1))

	for _, trustedPeer := range trustedList0 {
		require.True(t,
			lo.ContainsBy(trustedList1, func(tp apiclient.PeeringNodeIdentityResponse) bool {
				return tp.PeeringURL == trustedPeer.PeeringURL && tp.PublicKey == trustedPeer.PublicKey && tp.IsTrusted == trustedPeer.IsTrusted
			}),
		)
	}
}

func TestWaspCLICantPeerWithSelf(t *testing.T) {
	w := newWaspCLITest(t, waspClusterOpts{
		nNodes: 1,
	})

	peeringInfoOutput := w.MustRun("peering", "info")
	pubKey := regexp.MustCompile(`PubKey:\s+([[:alnum:]]+)$`).FindStringSubmatch(peeringInfoOutput[0])[1]

	require.Panics(
		t,
		func() {
			w.MustRun("peering", "trust", "self", pubKey, "0.0.0.0:4000")
		})
}

func TestWaspCLIListTrustDistrust(t *testing.T) {
	w := newWaspCLITest(t)
	out := w.MustRun("peering", "list-trusted", "--node=0")
	// one of the entries starts with "1", meaning node 0 trusts node 1
	containsNode1 := func(output []string) bool {
		for _, line := range output {
			if strings.HasPrefix(line, "1") {
				return true
			}
		}
		return false
	}
	require.True(t, containsNode1(out))

	// distrust node 1
	w.MustRun("peering", "distrust", "1", "--node=0")

	// 1 is not included anymore in the trusted list
	out = w.MustRun("peering", "list-trusted", "--node=0")
	// one of the entries starts with "1", meaning node 0 trusts node 1
	require.False(t, containsNode1(out))
}

func TestWaspCLIMintNativeToken(t *testing.T) {
	w := newWaspCLITest(t)

	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)
	w.MustRun("chain", "deposit", "base:100000000", "--node=0")

	out := w.MustRun(
		"chain", "create-native-token",
		"--max-supply=1000000",
		"--melted-tokens=0",
		"--minted-tokens=0",
		"--allowance=base:1000000",
		"--token-name=TEST",
		"--token-decimals=8",
		"--token-symbol=TS",
		"--node=0",
		"-o",
	)

	reqID := findRequestIDInOutput(out)
	require.NotEmpty(t, reqID)

	out = w.MustRun("chain", "request", reqID, "--node=0")
	require.Contains(t, strings.Join(out, "\n"), "Error: (empty)")
}

func TestWaspCLIRegisterERC20NativeTokenOnRemoteChain(t *testing.T) {
	w := newWaspCLITest(t)

	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)
	w.MustRun("chain", "deposit", "base:100000000", "--node=0")

	w.CreateL2NativeToken(&iotago.SimpleTokenScheme{
		MaximumSupply: big.NewInt(1000000),
		MeltedTokens:  big.NewInt(0),
		MintedTokens:  big.NewInt(0),
	}, "test", "test_symbol", 1)

	w.MustRun("chain", "deploy", "--chain=chain2", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain2", 0)
	w.MustRun("chain", "deposit", "base:100000000", "--node=0", "--chain=chain2")

	out := w.MustRun(
		"chain", "register-erc20-native-token-on-remote-chain",
		"-o",
		"--foundry-sn=1",
		"--token-name=test",
		"--ticker-symbol=test_symbol",
		"--token-decimals=1",
		"--target=chain2",
		"--node=0",
		"--chain=chain1",
		"--allowance=base:1000000",
	)

	reqID := findRequestIDInOutput(out)
	require.NotEmpty(t, reqID)

	out = w.MustRun("chain", "request", reqID, "--node=0", "--chain=chain1")
	require.Contains(t, strings.Join(out, "\n"), "Error: (empty)")
}

func sendDummyEVMTx(t *testing.T, w *WaspCLITest, ethPvtKey *ecdsa.PrivateKey) *types.Transaction {
	gasPrice := gas.DefaultFeePolicy().DefaultGasPriceFullDecimals(parameters.L1().BaseToken.Decimals)
	jsonRPCClient := NewEVMJSONRPClient(t, w.ChainID(0), w.Cluster, 0)
	tx, err := types.SignTx(
		types.NewTransaction(0, common.Address{}, big.NewInt(123), 100000, gasPrice, []byte{}),
		EVMSigner(),
		ethPvtKey,
	)
	require.NoError(t, err)
	err = jsonRPCClient.SendTransaction(context.Background(), tx)
	require.NoError(t, err)
	return tx
}

func TestEVMISCReceipt(t *testing.T) {
	w := newWaspCLITest(t)
	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)
	ethPvtKey, ethAddr := newEthereumAccount()
	w.MustRun("chain", "deposit", ethAddr.String(), "base:100000000", "--node=0")

	// send some arbitrary EVM tx
	tx := sendDummyEVMTx(t, w, ethPvtKey)
	out := w.MustRun("chain", "request", tx.Hash().Hex(), "--node=0")
	require.Contains(t, out[0], "Request found in block")
}

func TestChangeGovernanceController(t *testing.T) {
	w := newWaspCLITest(t)
	committee, quorum := w.ArgCommitteeConfig(0)
	w.MustRun("chain", "deploy", "--chain=chain1", committee, quorum, "--node=0")
	w.ActivateChainOnAllNodes("chain1", 0)

	// create the new controller
	_, newGovControllerAddr := testkey.GenKeyAddr()
	// change gov controller
	w.MustRun("chain", "change-gov-controller", newGovControllerAddr.Bech32("atoi"), "--chain=chain1")

	outputs, err := w.Cluster.L1Client().OutputMap(newGovControllerAddr)
	require.NoError(t, err)
	require.Len(t, outputs, 1)
}
