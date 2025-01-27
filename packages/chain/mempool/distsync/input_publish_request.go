// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package distsync

import (
	"fmt"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
)

type inputPublishRequest struct {
	request isc.Request
}

func NewInputPublishRequest(request isc.Request) gpa.Input {
	return &inputPublishRequest{request: request}
}

func (inp *inputPublishRequest) String() string {
	return fmt.Sprintf("{distSync.inputPublishRequest, request.ID=%v}", inp.request.ID().String())
}
