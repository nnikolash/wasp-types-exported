package controllerutils

import (
	"github.com/labstack/echo/v4"

	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/apierrors"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/interfaces"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/params"
)

// both used by the prometheus metrics middleware
const (
	EchoContextKeyChainID   = "chainid"
	EchoContextKeyOperation = "operation"
)

func ChainIDFromParams(c echo.Context, cs interfaces.ChainService) (isc.ChainID, error) {
	chainID, err := params.DecodeChainID(c)
	if err != nil {
		return isc.ChainID{}, err
	}

	if !cs.HasChain(chainID) {
		return isc.ChainID{}, apierrors.ChainNotFoundError()
	}
	// set chainID to be used by the prometheus metrics
	c.Set(EchoContextKeyChainID, chainID)
	return chainID, nil
}

// sets the label of the operation (endpoint being called) to be used by the prometheus metrics middleware
func SetOperation(c echo.Context, op string) {
	c.Set(EchoContextKeyOperation, op)
}

func ChainFromParams(c echo.Context, cs interfaces.ChainService) (chain.Chain, isc.ChainID, error) {
	chainID, err := params.DecodeChainID(c)
	if err != nil {
		return nil, isc.ChainID{}, err
	}

	chain, err := cs.GetChainByID(chainID)
	if err != nil {
		return nil, isc.ChainID{}, err
	}

	// set chainID to be used by the prometheus metrics
	c.Set(EchoContextKeyChainID, chainID)
	return chain, chainID, nil
}
