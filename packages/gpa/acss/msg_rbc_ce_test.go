// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestMsgRBCCEPayloadSerialization(t *testing.T) {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	require.NoError(t, err)
	msg := &msgRBCCEPayload{
		gpa.BasicMessage{},
		nil,
		b,
		nil,
	}

	rwutil.ReadWriteTest(t, msg, new(msgRBCCEPayload))
}
