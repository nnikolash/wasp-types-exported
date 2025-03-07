package processors

import (
	"go.uber.org/dig"

	"github.com/iotaledger/hive.go/app"
	"github.com/nnikolash/wasp-types-exported/contracts/native/inccounter"
	"github.com/nnikolash/wasp-types-exported/packages/isc/coreutil"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/coreprocessors"
	"github.com/nnikolash/wasp-types-exported/packages/vm/processors"
)

func init() {
	Component = &app.Component{
		Name:    "Processors",
		Provide: provide,
	}
}

var Component *app.Component

func provide(c *dig.Container) error {
	type processorsConfigResult struct {
		dig.Out

		ProcessorsConfig *processors.Config
	}

	if err := c.Provide(func() processorsConfigResult {
		Component.LogInfo("Registering native contracts...")

		nativeContracts := []*coreutil.ContractProcessor{
			inccounter.Processor,
		}

		for _, c := range nativeContracts {
			Component.LogDebugf(
				"Registering native contract: name: '%s', program hash: %s\n",
				c.Contract.Name, c.Contract.ProgramHash.String(),
			)
		}

		return processorsConfigResult{
			ProcessorsConfig: coreprocessors.NewConfigWithCoreContracts().WithNativeContracts(nativeContracts...),
		}
	}); err != nil {
		Component.LogPanic(err)
	}

	return nil
}
