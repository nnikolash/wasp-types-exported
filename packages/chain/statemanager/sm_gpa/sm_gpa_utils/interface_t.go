package sm_gpa_utils

import (
	"github.com/nnikolash/wasp-types-exported/packages/state"
)

// May be used only for tests; deleting in production should not be available.
type TestBlockWAL interface {
	BlockWAL
	Delete(state.BlockHash) bool
}
