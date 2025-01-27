package solo

import (
	"testing"

	"github.com/nnikolash/wasp-types-exported/packages/vm/core/corecontracts"
)

func TestSoloBasic1(t *testing.T) {
	corecontracts.PrintWellKnownHnames()
	env := New(t, &InitOptions{Debug: true, PrintStackTrace: true})
	_ = env.NewChain()
}
