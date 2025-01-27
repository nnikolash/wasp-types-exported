package wallet

import (
	"github.com/spf13/cobra"

	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/wallet"
)

func initMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wallet-migrate (keychain)",
		Short: "Migrates a seed inside the config file to the keychain provider",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			wallet.Migrate(wallet.WalletProvider(args[0]))
		},
	}
}
