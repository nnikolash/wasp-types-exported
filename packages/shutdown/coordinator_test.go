package shutdown_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/shutdown"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testlogger"
)

func TestShutdownCoordinator(t *testing.T) {
	log := testlogger.NewLogger(t)
	defer log.Sync()
	sc := shutdown.NewCoordinator("test", log)
	require.True(t, sc.CheckNestedDone())

	sc1 := sc.Nested("1")
	require.False(t, sc.CheckNestedDone())
	sc1.Done()
	require.True(t, sc.CheckNestedDone())
	sc.WaitNested()
}
