package corecontracts

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/controllers/controllerutils"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/corecontracts"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/models"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/params"
)

func MapGovChainInfoResponse(chainInfo *isc.ChainInfo) models.GovChainInfoResponse {
	return models.GovChainInfoResponse{
		ChainID:      chainInfo.ChainID.String(),
		ChainOwnerID: chainInfo.ChainOwnerID.String(),
		GasFeePolicy: chainInfo.GasFeePolicy,
		GasLimits:    chainInfo.GasLimits,
		PublicURL:    chainInfo.PublicURL,
		Metadata: models.GovPublicChainMetadata{
			EVMJsonRPCURL:   chainInfo.Metadata.EVMJsonRPCURL,
			EVMWebSocketURL: chainInfo.Metadata.EVMWebSocketURL,
			Name:            chainInfo.Metadata.Name,
			Description:     chainInfo.Metadata.Description,
			Website:         chainInfo.Metadata.Website,
		},
	}
}

func (c *Controller) getChainInfo(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	chainInfo, err := corecontracts.GetChainInfo(ch, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	chainInfoResponse := MapGovChainInfoResponse(chainInfo)

	return e.JSON(http.StatusOK, chainInfoResponse)
}

func (c *Controller) getChainOwner(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	chainOwner, err := corecontracts.GetChainOwner(ch, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	chainOwnerResponse := models.GovChainOwnerResponse{
		ChainOwner: chainOwner.String(),
	}

	return e.JSON(http.StatusOK, chainOwnerResponse)
}

func (c *Controller) getAllowedStateControllerAddresses(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	addresses, err := corecontracts.GetAllowedStateControllerAddresses(ch, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	encodedAddresses := make([]string, len(addresses))

	for k, v := range addresses {
		encodedAddresses[k] = v.Bech32(parameters.L1().Protocol.Bech32HRP)
	}

	addressesResponse := models.GovAllowedStateControllerAddressesResponse{
		Addresses: encodedAddresses,
	}

	return e.JSON(http.StatusOK, addressesResponse)
}
