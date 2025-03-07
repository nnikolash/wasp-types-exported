package sm_messages

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/chain/statemanager/sm_gpa/sm_gpa_utils"
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestBlockMessageSerialization(t *testing.T) {
	blocks := sm_gpa_utils.NewBlockFactory(t).GetBlocks(4, 1)
	for i := range blocks {
		// note that sender/receiver node IDs are transient
		// so don't use a random non-null node id here
		msg := NewBlockMessage(blocks[i], gpa.NodeID{})
		rwutil.ReadWriteTest(t, msg, NewEmptyBlockMessage())
	}
}

func TestSerializationBlockMessage(t *testing.T) {
	msg := &BlockMessage{
		gpa.BasicMessage{},
		state.RandomBlock(),
	}

	rwutil.ReadWriteTest(t, msg, new(BlockMessage))
}
