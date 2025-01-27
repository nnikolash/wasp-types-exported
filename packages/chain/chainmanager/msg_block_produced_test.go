package chainmanager

import (
	"testing"

	"github.com/iotaledger/iota.go/v3/tpkg"
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestMsgBlockProducedSerialization(t *testing.T) {
	msg := &msgBlockProduced{
		gpa.BasicMessage{},
		tpkg.RandTransaction(),
		state.RandomBlock(),
	}

	rwutil.ReadWriteTest(t, msg, new(msgBlockProduced))
}
