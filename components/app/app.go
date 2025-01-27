package app

import (
	_ "net/http/pprof"

	"github.com/iotaledger/hive.go/app"
	"github.com/iotaledger/hive.go/app/components/profiling"
	"github.com/iotaledger/hive.go/app/components/shutdown"
	"github.com/nnikolash/wasp-types-exported/components/cache"
	"github.com/nnikolash/wasp-types-exported/components/chains"
	"github.com/nnikolash/wasp-types-exported/components/database"
	"github.com/nnikolash/wasp-types-exported/components/dkg"
	"github.com/nnikolash/wasp-types-exported/components/logger"
	"github.com/nnikolash/wasp-types-exported/components/nodeconn"
	"github.com/nnikolash/wasp-types-exported/components/peering"
	"github.com/nnikolash/wasp-types-exported/components/processors"
	"github.com/nnikolash/wasp-types-exported/components/profilingrecorder"
	"github.com/nnikolash/wasp-types-exported/components/prometheus"
	"github.com/nnikolash/wasp-types-exported/components/publisher"
	"github.com/nnikolash/wasp-types-exported/components/registry"
	"github.com/nnikolash/wasp-types-exported/components/users"
	"github.com/nnikolash/wasp-types-exported/components/wasmtimevm"
	"github.com/nnikolash/wasp-types-exported/components/webapi"
	"github.com/nnikolash/wasp-types-exported/packages/toolset"
)

var (
	// Name of the app.
	Name = "Wasp"

	// Version of the app.
	// This field is populated by the scripts that compile wasp.
	Version = ""
)

func App() *app.App {
	return app.New(Name, Version,
		app.WithVersionCheck("iotaledger", "wasp"),
		app.WithInitComponent(InitComponent),
		app.WithComponents(
			shutdown.Component,
			nodeconn.Component,
			users.Component,
			logger.Component,
			cache.Component,
			database.Component,
			registry.Component,
			peering.Component,
			dkg.Component,
			processors.Component,
			wasmtimevm.Component,
			chains.Component,
			publisher.Component,
			webapi.Component,
			profiling.Component,
			profilingrecorder.Component,
			prometheus.Component,
		),
	)
}

var InitComponent *app.InitComponent

func init() {
	InitComponent = &app.InitComponent{
		Component: &app.Component{
			Name: "App",
		},
		NonHiddenFlags: []string{
			"app.checkForUpdates",
			"app.profile",
			"config",
			"help",
			"peering",
			"version",
		},
		AdditionalConfigs: []*app.ConfigurationSet{
			app.NewConfigurationSet("users", "users", "usersConfigFilePath", "usersConfig", false, true, false, "users.json", "u"),
		},
		Init: initialize,
	}
}

func initialize(_ *app.App) error {
	if toolset.ShouldHandleTools() {
		toolset.HandleTools()
		// HandleTools will call os.Exit
	}

	return nil
}
