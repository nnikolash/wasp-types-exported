// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"context"
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/packages/apilib"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/util/byz_quorum"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/cliclients"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/waspcmd"
)

func initRunDKGCmd() *cobra.Command {
	var (
		node   string
		peers  []string
		quorum int
	)

	cmd := &cobra.Command{
		Use:   "rundkg --peers=...",
		Short: "Runs the DKG on specified nodes",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			doDKG(node, peers, quorum)
		},
	}

	waspcmd.WithWaspNodeFlag(cmd, &node)
	waspcmd.WithPeersFlag(cmd, &peers)
	log.Check(cmd.MarkFlagRequired("peers"))
	cmd.Flags().IntVarP(&quorum, "quorum", "", 0, "quorum (default: 2/3s of the number of committee nodes)")
	return cmd
}

func doDKG(node string, peers []string, quorum int) iotago.Address {
	client := cliclients.WaspClient(node)
	nodeInfo, _, err := client.NodeApi.GetPeeringIdentity(context.Background()).Execute() //nolint:bodyclose // false positive
	log.Check(err)

	// Consider own node as a committee, if peers are not specified.
	if len(peers) == 0 {
		peers = append(peers, nodeInfo.PublicKey)
	}

	// grab the peering info of the peers from the node
	filteredPeers := make([]apiclient.PeeringNodeIdentityResponse, 0)
	thisNodeFound := false
	{
		var trustedPeers []apiclient.PeeringNodeIdentityResponse
		trustedPeers, _, err = client.NodeApi.GetTrustedPeers(context.Background()).Execute() //nolint:bodyclose // false positive
		log.Check(err)

		for _, peer := range peers {
			foundPeer, exists := lo.Find(trustedPeers, func(p apiclient.PeeringNodeIdentityResponse) bool {
				return (p.Name == peer || p.PublicKey == peer) && p.IsTrusted
			})
			if !exists {
				log.Fatalf("peer with name {%s} not found in trusted peers", peer)
			}
			if foundPeer.PublicKey == nodeInfo.PublicKey {
				thisNodeFound = true
			}
			filteredPeers = append(filteredPeers, foundPeer)
		}
	}
	if !thisNodeFound {
		// TODO: This is temporary, until DKG is fixed to not require the current node in the committee.
		fmt.Fprintf(os.Stdout, "NOTE: Adding this node as a committee member.\n")
		filteredPeers = append(filteredPeers, *nodeInfo)
	}

	committeePubKeys := []string{}
	for _, peer := range filteredPeers {
		committeePubKeys = append(committeePubKeys, peer.PublicKey)
	}

	// Use default quorum, if it is unspecified.
	minQuorum := byz_quorum.MinQuorum(len(committeePubKeys))
	if quorum == 0 {
		quorum = minQuorum
	}

	if quorum < minQuorum {
		log.Fatal("quorum needs to be at least (2/3)+1 of committee size")
	}

	stateControllerAddr, err := apilib.RunDKG(client, committeePubKeys, uint16(quorum))
	log.Check(err)

	committeeMembersStr := ""
	for _, fp := range filteredPeers {
		committeeMembersStr += fmt.Sprintf("%v (%v)\n", fp.PublicKey, fp.Name)
	}

	fmt.Fprintf(os.Stdout,
		"DKG successful\nAddress: %s\n* committee size = %v\n* quorum = %v\n* members: %s\n",
		stateControllerAddr.Bech32(parameters.L1().Protocol.Bech32HRP),
		len(committeePubKeys),
		quorum,
		committeeMembersStr,
	)
	return stateControllerAddr
}
