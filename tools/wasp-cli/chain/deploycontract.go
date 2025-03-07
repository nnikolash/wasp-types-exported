package chain

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/clients/chainclient"
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blob"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/root"
	"github.com/nnikolash/wasp-types-exported/packages/vm/vmtypes"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/cliclients"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/config"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/util"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/waspcmd"
)

func initDeployContractCmd() *cobra.Command {
	var node string
	var chain string

	cmd := &cobra.Command{
		Use:   "deploy-contract <vmtype> <name> <description> <filename|program-hash> [init-params]",
		Short: "Deploy a contract in the chain",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			chainID := config.GetChain(chain)
			client := cliclients.WaspClient(node)
			vmtype := args[0]
			name := args[1]
			description := args[2]
			initParams := util.EncodeParams(args[4:], chainID)

			var progHash hashing.HashValue

			switch vmtype {
			case vmtypes.Core:
				log.Fatal("cannot manually deploy core contracts")

			case vmtypes.Native:
				var err error
				progHash, err = hashing.HashValueFromHex(args[3])
				log.Check(err)

			default:
				filename := args[3]
				blobFieldValues := codec.MakeDict(map[string]interface{}{
					blob.VarFieldVMType:             vmtype,
					blob.VarFieldProgramDescription: description,
					blob.VarFieldProgramBinary:      util.ReadFile(filename),
				})
				progHash = uploadBlob(client, chainID, blobFieldValues)
			}
			deployContract(client, chainID, node, name, progHash, initParams)
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func deployContract(client *apiclient.APIClient, chainID isc.ChainID, node, name string, progHash hashing.HashValue, initParams dict.Dict) {
	util.WithOffLedgerRequest(chainID, node, func() (isc.OffLedgerRequest, error) {
		args := codec.MakeDict(map[string]interface{}{
			root.ParamName:        name,
			root.ParamProgramHash: progHash,
		})
		args.Extend(initParams)
		return cliclients.ChainClient(client, chainID).PostOffLedgerRequest(context.Background(),
			root.Contract.Hname(),
			root.FuncDeployContract.Hname(),
			chainclient.PostRequestParams{
				Args: args,
			},
		)
	})
}
