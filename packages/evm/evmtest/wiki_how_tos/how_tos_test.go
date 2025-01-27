package wiki_how_tos_test

import (
	_ "embed"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/solo"
	"github.com/nnikolash/wasp-types-exported/packages/util"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/evmtest"
)

//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` GetBalance.sol -o ."
var (
	//go:embed GetBalance.abi
	GetBalanceContractABI string
	//go:embed GetBalance.bin
	GetBalanceContractBytecodeHex string
	GetBalanceContractBytecode    = common.FromHex(strings.TrimSpace(GetBalanceContractBytecodeHex))
)

//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` Entropy.sol -o ."
var (
	//go:embed Entropy.abi
	EntropyContractABI string
	//go:embed Entropy.bin
	EntropyContractBytecodeHex string
	EntropyContractBytecode    = common.FromHex(strings.TrimSpace(EntropyContractBytecodeHex))
)

func TestBaseBalance(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	balance, _ := env.Chain.EVM().Balance(deployer, nil)
	decimals := env.Chain.EVM().BaseToken().Decimals
	var value uint64
	instance.CallFnExpectEvent(nil, "GotBaseBalance", &value, "getBalanceBaseTokens")
	realBalance := util.BaseTokensDecimalsToEthereumDecimals(value, decimals)
	assert.Equal(t, balance, realBalance)
}

func TestNativeBalance(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	// create a new native token on L1
	foundry, tokenID, err := env.Chain.NewNativeTokenParams(100000000000000).CreateFoundry()
	require.NoError(t, err)
	// the token id in bytes, used to call the contract
	nativeTokenIDBytes := isc.NativeTokenIDToBytes(tokenID)

	// mint some native tokens to the chain originator
	err = env.Chain.MintTokens(foundry, 10000000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// get the agentId of the contract deployer
	senderAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// send some native tokens to the contract deployer
	// and check if the balance returned by the contract is correct
	err = env.Chain.SendFromL2ToL2AccountNativeTokens(tokenID, senderAgentID, 100000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	nativeBalance := new(big.Int)
	instance.CallFnExpectEvent(nil, "GotNativeTokenBalance", &nativeBalance, "getBalanceNativeTokens", nativeTokenIDBytes)
	assert.Equal(t, int64(100000), nativeBalance.Int64())
}

func TestNFTBalance(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	// get the agentId of the contract deployer
	senderAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// mint an NFToken to the contract deployer
	// and check if the balance returned by the contract is correct
	mockMetaData := []byte("sesa")
	nfti, info, err := env.Chain.Env.MintNFTL1(env.Chain.OriginatorPrivateKey, env.Chain.OriginatorAddress, mockMetaData)
	require.NoError(t, err)
	env.Chain.MustDepositNFT(nfti, env.Chain.OriginatorAgentID, env.Chain.OriginatorPrivateKey)

	transfer := isc.NewEmptyAssets()
	transfer.AddNFTs(info.NFTID)

	// send the NFT to the contract deployer
	err = env.Chain.SendFromL2ToL2Account(transfer, senderAgentID, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// get the NFT balance of the contract deployer
	nftBalance := new(big.Int)
	instance.CallFnExpectEvent(nil, "GotNFTIDs", &nftBalance, "getBalanceNFTs")
	assert.Equal(t, int64(1), nftBalance.Int64())
}

func TestAgentID(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetBalanceContractABI, GetBalanceContractBytecode)

	// get the agentId of the contract deployer
	senderAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)

	// get the agnetId of the contract deployer
	// and compare it with the agentId returned by the contract
	var agentID []byte
	instance.CallFnExpectEvent(nil, "GotAgentID", &agentID, "getAgentID")
	assert.Equal(t, senderAgentID.Bytes(), agentID)
}

func TestEntropy(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, _ := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, EntropyContractABI, EntropyContractBytecode)

	// get the entropy of the contract
	// and check if it is different from the previous one
	var entropy [32]byte
	instance.CallFnExpectEvent(nil, "EntropyEvent", &entropy, "emitEntropy")
	var entropy2 [32]byte
	instance.CallFnExpectEvent(nil, "EntropyEvent", &entropy2, "emitEntropy")
	assert.NotEqual(t, entropy, entropy2)
}
