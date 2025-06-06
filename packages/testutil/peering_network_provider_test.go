// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package testutil_test

import (
	"testing"
	"time"

	"github.com/nnikolash/wasp-types-exported/packages/peering"
	"github.com/nnikolash/wasp-types-exported/packages/testutil"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testlogger"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testpeers"
)

func TestFakeNetwork(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()
	doneCh := make(chan bool)
	chain1 := peering.RandomPeeringID()
	chain2 := peering.RandomPeeringID()
	receiver := byte(0)
	peeringURLs, nodeIdentities := testpeers.SetupKeys(3)
	network := testutil.NewPeeringNetwork(peeringURLs, nodeIdentities, 100, testutil.NewPeeringNetReliable(log), log)
	netProviders := network.NetworkProviders()
	//
	// Node "a" listens for chain1 messages.
	netProviders[0].Attach(&chain1, receiver, func(recv *peering.PeerMessageIn) {
		doneCh <- true
	})
	//
	// Node "b" sends some messages.
	var a, c peering.PeerSender
	a, _ = netProviders[1].PeerByPubKey(nodeIdentities[0].GetPublicKey())
	c, _ = netProviders[1].PeerByPubKey(nodeIdentities[2].GetPublicKey())
	a.SendMsg(peering.NewPeerMessageData(chain1, receiver, 1)) // Will be delivered.
	a.SendMsg(peering.NewPeerMessageData(chain2, receiver, 2)) // Will be dropped.
	a.SendMsg(peering.NewPeerMessageData(chain1, byte(5), 3))  // Will be dropped.
	c.SendMsg(peering.NewPeerMessageData(chain1, receiver, 4)) // Will be dropped.
	//
	// Wait for the result.
	select {
	case <-doneCh:
	case <-time.After(1 * time.Second):
		panic("timeout")
	}
}
