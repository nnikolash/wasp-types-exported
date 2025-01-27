// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestMsgVoteSerialization(t *testing.T) {
	{
		msg := &msgVote{
			gpa.BasicMessage{},
			msgVoteOK,
		}

		rwutil.ReadWriteTest(t, msg, new(msgVote))
	}
	{
		msg := &msgVote{
			gpa.BasicMessage{},
			msgVoteREADY,
		}

		rwutil.ReadWriteTest(t, msg, new(msgVote))
	}
}
