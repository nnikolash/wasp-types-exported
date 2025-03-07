package chain

import (
	"encoding/hex"

	"github.com/spf13/cobra"

	"github.com/nnikolash/wasp-types-exported/clients/chainclient"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/config"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/util"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/waspcmd"
)

func initRegisterERC20NativeTokenCmd() *cobra.Command {
	return buildPostRequestCmd(
		"register-erc20-native-token",
		"Call evm core contract registerERC20NativeToken entry point",
		evm.Contract.Name,
		evm.FuncRegisterERC20NativeToken.Name,
		func(cmd *cobra.Command) {
			initRegisterERC20NativeTokenParams(cmd)
		},
		getRegisterERC20NativeTokenArgs,
	)
}

func initRegisterERC20NativeTokenOnRemoteChainCmd() *cobra.Command {
	var targetChain string

	return buildPostRequestCmd(
		"register-erc20-native-token-on-remote-chain",
		"Call evm core contract registerERC20NativeTokenOnRemoteChain entry point",
		evm.Contract.Name,
		evm.FuncRegisterERC20NativeTokenOnRemoteChain.Name,
		func(cmd *cobra.Command) {
			initRegisterERC20NativeTokenParams(cmd)
			cmd.Flags().StringVarP(&targetChain, "target", "A", "", "Target chain ID")
		},
		func(cmd *cobra.Command) []string {
			chainID := codec.EncodeAddress(config.GetChain(targetChain).AsAddress())
			extraArgs := []string{"string", "A", "bytes", "0x" + hex.EncodeToString(chainID)}
			return append(getRegisterERC20NativeTokenArgs(cmd), extraArgs...)
		},
	)
}

func buildPostRequestCmd(name, desc, hname, fname string, initFlags func(cmd *cobra.Command), funcArgs func(cmd *cobra.Command) []string) *cobra.Command {
	var (
		chain             string
		node              string
		postrequestParams postRequestParams
	)

	cmd := &cobra.Command{
		Use:   name,
		Short: desc,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)
			chainID := config.GetChain(chain)

			allowanceTokens := util.ParseFungibleTokens(postrequestParams.allowance)

			params := chainclient.PostRequestParams{
				Args:      util.EncodeParams(funcArgs(cmd), chainID),
				Transfer:  util.ParseFungibleTokens(postrequestParams.transfer),
				Allowance: allowanceTokens,
			}
			postRequest(
				node,
				chain,
				hname,
				fname,
				params,
				postrequestParams.offLedger,
				postrequestParams.adjustStorageDeposit,
			)
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	postrequestParams.initFlags(cmd)
	initFlags(cmd)

	return cmd
}

func initRegisterERC20NativeTokenParams(cmd *cobra.Command) {
	cmd.Flags().Uint32("foundry-sn", 0, "Foundry serial number")
	cmd.Flags().String("token-name", "", "Token name")
	cmd.Flags().String("ticker-symbol", "", "Ticker symbol")
	cmd.Flags().Uint8("token-decimals", 0, "Token decimals")
}

func getRegisterERC20NativeTokenArgs(cmd *cobra.Command) []string {
	return []string{
		"string", "fs", "uint32", flagValString(cmd, "foundry-sn"),
		"string", "n", "string", flagValString(cmd, "token-name"),
		"string", "t", "string", flagValString(cmd, "ticker-symbol"),
		"string", "d", "uint8", flagValString(cmd, "token-decimals"),
	}
}

func flagValString(cmd *cobra.Command, cmdName string) string {
	return cmd.Flag(cmdName).Value.String()
}
