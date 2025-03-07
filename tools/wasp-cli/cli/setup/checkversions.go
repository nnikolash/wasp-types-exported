package setup

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nnikolash/wasp-types-exported/components/app"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/cliclients"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
)

func initCheckVersionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-versions",
		Short: "checks the versions of wasp-cli and wasp nodes match",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// query every wasp node info endpoint and ensure the `Version` matches
			waspSettings := map[string]interface{}{}
			waspKey := viper.Sub("wasp")
			if waspKey != nil {
				waspSettings = waspKey.AllSettings()
			}
			if len(waspSettings) == 0 {
				log.Fatalf("no wasp node configured, you can add a node with `wasp-cli wasp add <name> <api url>`")
			}
			for nodeName := range waspSettings {
				nodeVersion, _, err := cliclients.WaspClient(nodeName).NodeApi.
					GetVersion(context.Background()).
					Execute()
				log.Check(err)
				if app.Version == "v"+nodeVersion.Version {
					log.Printf("Wasp-cli version matches Wasp {%s}\n", nodeName)
				} else {
					log.Printf("! -> Version mismatch with Wasp {%s}. cli version: %s, wasp version: %s\n", nodeName, app.Version, nodeVersion.Version)
				}
			}
		},
	}
}
