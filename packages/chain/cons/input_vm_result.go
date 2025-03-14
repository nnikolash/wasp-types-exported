// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"fmt"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
	"github.com/nnikolash/wasp-types-exported/packages/vm"
)

type inputVMResult struct {
	task *vm.VMTaskResult
}

func NewInputVMResult(task *vm.VMTaskResult) gpa.Input {
	return &inputVMResult{task: task}
}

func (inp *inputVMResult) String() string {
	return fmt.Sprintf("{cons.inputVMResult: %+v}", inp.task)
}
