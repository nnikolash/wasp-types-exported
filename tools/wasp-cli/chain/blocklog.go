package chain

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/cliclients"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/config"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/util"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/waspcmd"
)

func initBlockCmd() *cobra.Command {
	var node string
	var chain string
	cmd := &cobra.Command{
		Use:   "block [index]",
		Short: "Get information about a block given its index, or latest block if missing",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			bi := fetchBlockInfo(args, node, chain)
			log.Printf("Block index: %d\n", bi.BlockIndex)
			log.Printf("Timestamp: %s\n", bi.Timestamp.UTC().Format(time.RFC3339))
			log.Printf("Total requests: %d\n", bi.TotalRequests)
			log.Printf("Successful requests: %d\n", bi.NumSuccessfulRequests)
			log.Printf("Off-ledger requests: %d\n", bi.NumOffLedgerRequests)
			log.Printf("\n")
			logRequestsInBlock(bi.BlockIndex, node, chain)
			log.Printf("\n")
			logEventsInBlock(bi.BlockIndex, node, chain)
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func fetchBlockInfo(args []string, node, chain string) *apiclient.BlockInfoResponse {
	client := cliclients.WaspClient(node)

	if len(args) == 0 {
		blockInfo, _, err := client.
			CorecontractsApi.
			BlocklogGetLatestBlockInfo(context.Background(), config.GetChain(chain).String()).
			Execute() //nolint:bodyclose // false positive

		log.Check(err)
		return blockInfo
	}

	blockIndexStr := args[0]
	index, err := strconv.ParseUint(blockIndexStr, 10, 32)
	log.Check(err)

	blockInfo, _, err := client.
		CorecontractsApi.
		BlocklogGetBlockInfo(context.Background(), config.GetChain(chain).String(), uint32(index)).
		Block(blockIndexStr).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)
	return blockInfo
}

func logRequestsInBlock(index uint32, node, chain string) {
	client := cliclients.WaspClient(node)
	receipts, _, err := client.CorecontractsApi.
		BlocklogGetRequestReceiptsOfBlock(context.Background(), config.GetChain(chain).String(), index).
		Block(fmt.Sprintf("%d", index)).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)

	for i, receipt := range receipts {
		r := receipt
		util.LogReceipt(r, i)
	}
}

func logEventsInBlock(index uint32, node, chain string) {
	client := cliclients.WaspClient(node)
	events, _, err := client.CorecontractsApi.
		BlocklogGetEventsOfBlock(context.Background(), config.GetChain(chain).String(), index).
		Block(fmt.Sprintf("%d", index)).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)
	logEvents(events)
}

func hexLenFromByteLen(length int) int {
	return (length * 2) + 2
}

func reqIDFromString(s string) isc.RequestID {
	switch len(s) {
	case hexLenFromByteLen(iotago.OutputIDLength):
		// isc ReqID
		reqID, err := isc.RequestIDFromString(s)
		log.Check(err)
		return reqID
	case hexLenFromByteLen(common.HashLength):
		bytes, err := iotago.DecodeHex(s)
		log.Check(err)
		var txHash common.Hash
		copy(txHash[:], bytes)
		return isc.RequestIDFromEVMTxHash(txHash)
	default:
		log.Fatalf("invalid requestID length: %d", len(s))
	}
	panic("unreachable")
}

func initRequestCmd() *cobra.Command {
	var node string
	var chain string
	cmd := &cobra.Command{
		Use:   "request <request-id>",
		Short: "Get information about a request given its ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)
			chainID := config.GetChain(chain)

			client := cliclients.WaspClient(node)
			reqID := reqIDFromString(args[0])

			// TODO add optional block param?
			receipt, _, err := client.ChainsApi.
				GetReceipt(context.Background(), chainID.String(), reqID.String()).
				Execute() //nolint:bodyclose // false positive

			log.Check(err)

			log.Printf("Request found in block %d\n\n", receipt.BlockIndex)
			util.LogReceipt(*receipt)

			log.Printf("\n")
			logEventsInRequest(reqID, node, chain)
			log.Printf("\n")
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func logEventsInRequest(reqID isc.RequestID, node, chain string) {
	client := cliclients.WaspClient(node)
	events, _, err := client.CorecontractsApi.
		BlocklogGetEventsOfRequest(context.Background(), config.GetChain(chain).String(), reqID.String()).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)
	logEvents(events)
}

func logEvents(ret *apiclient.EventsResponse) {
	header := []string{"event"}
	rows := make([][]string, len(ret.Events))

	for i, event := range ret.Events {
		rows[i] = []string{event.Topic}
	}

	log.Printf("Total %d events\n", len(ret.Events))
	log.PrintTable(header, rows)
}
