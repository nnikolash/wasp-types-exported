package cliclients

import (
	"context"

	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/clients/apiextensions"
	"github.com/nnikolash/wasp-types-exported/clients/chainclient"
	"github.com/nnikolash/wasp-types-exported/clients/scclient"
	"github.com/nnikolash/wasp-types-exported/components/app"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/l1connection"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/config"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/wallet"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
)

var SkipCheckVersions bool

func WaspClientForHostName(name string) *apiclient.APIClient {
	apiAddress := config.MustWaspAPIURL(name)
	L1Client() // this will fill parameters.L1() with data from the L1 node
	log.Verbosef("using Wasp host %s\n", apiAddress)

	client, err := apiextensions.WaspAPIClientByHostName(apiAddress)
	log.Check(err)

	client.GetConfig().Debug = log.DebugFlag
	client.GetConfig().AddDefaultHeader("Authorization", "Bearer "+config.GetToken(name))

	return client
}

func WaspClient(name string) *apiclient.APIClient {
	client := WaspClientForHostName(name)
	assertMatchingNodeVersion(name, client)
	return client
}

func assertMatchingNodeVersion(name string, client *apiclient.APIClient) {
	if SkipCheckVersions {
		return
	}
	nodeVersion, _, err := client.NodeApi.
		GetVersion(context.Background()).
		Execute()
	log.Check(err)
	if app.Version != "v"+nodeVersion.Version {
		log.Fatalf("node [%s] version: %s, does not match wasp-cli version: %s. You can skip this check by re-running with command with --skip-version-check",
			name, nodeVersion.Version, app.Version)
	}
}

func L1Client() l1connection.Client {
	log.Verbosef("using L1 API %s\n", config.L1APIAddress())

	return l1connection.NewClient(
		l1connection.Config{
			APIAddress:    config.L1APIAddress(),
			FaucetAddress: config.L1FaucetAddress(),
		},
		log.HiveLogger(),
	)
}

func ChainClient(waspClient *apiclient.APIClient, chainID isc.ChainID) *chainclient.Client {
	return chainclient.New(
		L1Client(),
		waspClient,
		chainID,
		wallet.Load(),
	)
}

func SCClient(apiClient *apiclient.APIClient, chainID isc.ChainID, contractHname isc.Hname) *scclient.SCClient {
	return scclient.New(ChainClient(apiClient, chainID), contractHname)
}
