package chain

import (
	"encoding/hex"
	"math/big"

	"github.com/spf13/cobra"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
)

func initCreateNativeTokenCmd() *cobra.Command {
	var maxSupply, mintedTokens, meltedTokens int64
	var tokenName, tokenSymbol string
	var tokenDecimals uint8

	return buildPostRequestCmd(
		"create-native-token",
		"Calls accounts core contract nativeTokenCreate to create a new native token",
		accounts.Contract.Name,
		accounts.FuncNativeTokenCreate.Name,
		func(cmd *cobra.Command) {
			cmd.Flags().Int64Var(&maxSupply, "max-supply", 1000000, "Maximum token supply")
			cmd.Flags().Int64Var(&mintedTokens, "minted-tokens", 0, "Minted tokens")
			cmd.Flags().Int64Var(&meltedTokens, "melted-tokens", 0, "Melted tokens")
			cmd.Flags().StringVar(&tokenName, "token-name", "", "Token name")
			cmd.Flags().StringVar(&tokenSymbol, "token-symbol", "", "Token symbol")
			cmd.Flags().Uint8Var(&tokenDecimals, "token-decimals", uint8(8), "Token decimals")
		},
		func(cmd *cobra.Command) []string {
			tokenScheme := &iotago.SimpleTokenScheme{
				MaximumSupply: big.NewInt(maxSupply),
				MintedTokens:  big.NewInt(mintedTokens),
				MeltedTokens:  big.NewInt(meltedTokens),
			}

			tokenSchemeBytes := codec.EncodeTokenScheme(tokenScheme)

			return []string{
				"string", accounts.ParamTokenScheme, "bytes", "0x" + hex.EncodeToString(tokenSchemeBytes),
				"string", accounts.ParamTokenName, "bytes", "0x" + hex.EncodeToString(codec.EncodeString(tokenName)),
				"string", accounts.ParamTokenTickerSymbol, "bytes", "0x" + hex.EncodeToString(codec.EncodeString(tokenSymbol)),
				"string", accounts.ParamTokenDecimals, "bytes", "0x" + hex.EncodeToString(codec.EncodeUint8(tokenDecimals)),
			}
		},
	)
}
