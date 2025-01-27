package chain

import (
	"context"

	"github.com/spf13/cobra"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/cliclients"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/config"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/util"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/waspcmd"
)

var uploadQuorum int

func initUploadFlags(chainCmd *cobra.Command) {
	chainCmd.PersistentFlags().IntVarP(&uploadQuorum, "upload-quorum", "", 3, "quorum for blob upload")
}

func initStoreBlobCmd() *cobra.Command {
	var node string
	var chain string
	cmd := &cobra.Command{
		Use:   "store-blob <type> <field> <type> <value> ...",
		Short: "Store a blob in the chain",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			chainID := config.GetChain(chain)
			uploadBlob(cliclients.WaspClient(node), chainID, util.EncodeParams(args, chainID))
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func uploadBlob(client *apiclient.APIClient, chainID isc.ChainID, fieldValues dict.Dict) (hash hashing.HashValue) {
	chainClient := cliclients.ChainClient(client, chainID)

	hash, _, receipt, err := chainClient.UploadBlob(context.Background(), fieldValues)
	log.Check(err)
	log.Printf("uploaded blob to chain -- hash: %s\n", hash)
	util.LogReceipt(*receipt)
	return hash
}

func initShowBlobCmd() *cobra.Command {
	var node string
	var chain string
	cmd := &cobra.Command{
		Use:   "show-blob <hash>",
		Short: "Show a blob in chain",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			hash, err := hashing.HashValueFromHex(args[0])
			log.Check(err)

			client := cliclients.WaspClient(node)

			blobInfo, _, err := client.
				CorecontractsApi.
				BlobsGetBlobInfo(context.Background(), config.GetChain(chain).String(), hash.Hex()).
				Execute() //nolint:bodyclose // false positive
			log.Check(err)

			values := dict.New()
			for field := range blobInfo.Fields {
				value, _, err := client.
					CorecontractsApi.
					BlobsGetBlobValue(context.Background(), config.GetChain(chain).String(), hash.Hex(), field).
					Execute() //nolint:bodyclose // false positive

				log.Check(err)

				decodedValue, err := iotago.DecodeHex(value.ValueData)
				log.Check(err)

				values.Set(kv.Key(field), decodedValue)
			}
			util.PrintDictAsJSON(values)
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}
