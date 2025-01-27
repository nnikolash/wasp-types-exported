package wallet

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/config"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/wallet"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
)

type WalletConfig struct {
	KeyPair cryptolib.VariantKeyPair
}

var initOverwrite bool

func initInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new wallet",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			wallet.InitWallet(initOverwrite)

			config.SetWalletProviderString(string(wallet.GetWalletProvider()))
			log.Check(viper.WriteConfig())
		},
	}
	cmd.Flags().BoolVar(&initOverwrite, "overwrite", false, "allow overwriting existing seed")
	return cmd
}

type InitModel struct {
	Scheme string
}

var _ log.CLIOutput = &InitModel{}

func (i *InitModel) AsText() (string, error) {
	template := `Initialized wallet seed in {{ .Scheme }}`
	return log.ParseCLIOutputTemplate(i, template)
}
