package peering_test

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/peering"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestPeerMessageSerialization(t *testing.T) {
	msg := &peering.PeerMessageNet{
		PeerMessageData: peering.NewPeerMessageData(
			peering.RandomPeeringID(),
			byte(10),
			peering.FirstUserMsgCode+17,
			[]byte{1, 2, 3, 4, 5}),
	}
	rwutil.BytesTest(t, msg, peering.PeerMessageNetFromBytes)
}
