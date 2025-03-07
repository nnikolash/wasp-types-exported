package chain

import (
	"github.com/spf13/cobra"

	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
)

func initChainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "chain <command>",
		Short: "Interact with a chain",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log.Check(cmd.Help())
		},
	}
}

func Init(rootCmd *cobra.Command) {
	chainCmd := initChainCmd()
	rootCmd.AddCommand(chainCmd)

	initUploadFlags(chainCmd)

	chainCmd.AddCommand(initListCmd())
	chainCmd.AddCommand(initDeployCmd())
	chainCmd.AddCommand(initInfoCmd())
	chainCmd.AddCommand(initListContractsCmd())
	chainCmd.AddCommand(initDeployContractCmd())
	chainCmd.AddCommand(initBalanceCmd())
	chainCmd.AddCommand(initAccountNFTsCmd())
	chainCmd.AddCommand(initDepositCmd())
	chainCmd.AddCommand(initStoreBlobCmd())
	chainCmd.AddCommand(initShowBlobCmd())
	chainCmd.AddCommand(initBlockCmd())
	chainCmd.AddCommand(initRequestCmd())
	chainCmd.AddCommand(initPostRequestCmd())
	chainCmd.AddCommand(initCallViewCmd())
	chainCmd.AddCommand(initActivateCmd())
	chainCmd.AddCommand(initDeactivateCmd())
	chainCmd.AddCommand(initRunDKGCmd())
	chainCmd.AddCommand(initRotateCmd())
	chainCmd.AddCommand(initRotateWithDKGCmd())
	chainCmd.AddCommand(initChangeGovControllerCmd())
	chainCmd.AddCommand(initChangeAccessNodesCmd())
	chainCmd.AddCommand(initDisableFeePolicyCmd())
	chainCmd.AddCommand(initPermissionlessAccessNodesCmd())
	chainCmd.AddCommand(initAddChainCmd())
	chainCmd.AddCommand(initRegisterERC20NativeTokenCmd())
	chainCmd.AddCommand(initRegisterERC20NativeTokenOnRemoteChainCmd())
	chainCmd.AddCommand(initCreateNativeTokenCmd())
	chainCmd.AddCommand(initMetadataCmd())
}
