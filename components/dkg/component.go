// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package dkg

import (
	"fmt"

	"go.uber.org/dig"

	"github.com/iotaledger/hive.go/app"
	"github.com/nnikolash/wasp-types-exported/packages/dkg"
	"github.com/nnikolash/wasp-types-exported/packages/peering"
	"github.com/nnikolash/wasp-types-exported/packages/registry"
)

func init() {
	Component = &app.Component{
		Name:    "DKG",
		Provide: provide,
	}
}

var Component *app.Component

func provide(c *dig.Container) error {
	type nodeDeps struct {
		dig.In

		NodeIdentityProvider    registry.NodeIdentityProvider
		DKShareRegistryProvider registry.DKShareRegistryProvider
		NetworkProvider         peering.NetworkProvider `name:"networkProvider"`
	}

	type nodeResult struct {
		dig.Out

		Node *dkg.Node
	}

	if err := c.Provide(func(deps nodeDeps) nodeResult {
		node, err := dkg.NewNode(
			deps.NodeIdentityProvider.NodeIdentity(),
			deps.NetworkProvider,
			deps.DKShareRegistryProvider,
			Component.Logger(),
		)
		if err != nil {
			Component.LogPanic(fmt.Errorf("failed to initialize the DKG node: %w", err))
		}

		return nodeResult{
			Node: node,
		}
	}); err != nil {
		Component.LogPanic(err)
	}

	return nil
}
