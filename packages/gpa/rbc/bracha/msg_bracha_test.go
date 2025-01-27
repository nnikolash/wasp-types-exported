// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package bracha

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestMsgBrachaSerialization(t *testing.T) {
	{
		b := make([]byte, 10)
		_, err := rand.Read(b)
		require.NoError(t, err)
		msg := &msgBracha{
			gpa.BasicMessage{},
			msgBrachaTypePropose,
			b,
		}

		rwutil.ReadWriteTest(t, msg, new(msgBracha))
	}
	{
		b := make([]byte, 10)
		_, err := rand.Read(b)
		require.NoError(t, err)
		msg := &msgBracha{
			gpa.BasicMessage{},
			msgBrachaTypeEcho,
			b,
		}

		rwutil.ReadWriteTest(t, msg, new(msgBracha))
	}
	{
		b := make([]byte, 10)
		_, err := rand.Read(b)
		require.NoError(t, err)
		msg := &msgBracha{
			gpa.BasicMessage{},
			msgBrachaTypeReady,
			b,
		}

		rwutil.ReadWriteTest(t, msg, new(msgBracha))
	}
}
