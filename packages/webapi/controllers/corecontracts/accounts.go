package corecontracts

import (
	"net/http"

	"github.com/labstack/echo/v4"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/apierrors"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/controllers/controllerutils"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/corecontracts"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/models"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/params"
)

func (c *Controller) getTotalAssets(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	assets, err := corecontracts.GetTotalAssets(ch, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	assetsResponse := &models.AssetsResponse{
		BaseTokens:   iotago.EncodeUint64(assets.BaseTokens),
		NativeTokens: isc.NativeTokensToJSONObject(assets.NativeTokens),
	}

	return e.JSON(http.StatusOK, assetsResponse)
}

func (c *Controller) getAccountBalance(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	agentID, err := params.DecodeAgentID(e)
	if err != nil {
		return err
	}

	assets, err := corecontracts.GetAccountBalance(ch, agentID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	assetsResponse := &models.AssetsResponse{
		BaseTokens:   iotago.EncodeUint64(assets.BaseTokens),
		NativeTokens: isc.NativeTokensToJSONObject(assets.NativeTokens),
	}

	return e.JSON(http.StatusOK, assetsResponse)
}

func (c *Controller) getAccountNFTs(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	agentID, err := params.DecodeAgentID(e)
	if err != nil {
		return err
	}

	nfts, err := corecontracts.GetAccountNFTs(ch, agentID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	nftsResponse := &models.AccountNFTsResponse{
		NFTIDs: make([]string, len(nfts)),
	}

	for k, v := range nfts {
		nftsResponse.NFTIDs[k] = v.ToHex()
	}

	return e.JSON(http.StatusOK, nftsResponse)
}

func (c *Controller) getAccountFoundries(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}
	agentID, err := params.DecodeAgentID(e)
	if err != nil {
		return err
	}

	foundries, err := corecontracts.GetAccountFoundries(ch, agentID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	return e.JSON(http.StatusOK, &models.AccountFoundriesResponse{
		FoundrySerialNumbers: foundries,
	})
}

func (c *Controller) getAccountNonce(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	agentID, err := params.DecodeAgentID(e)
	if err != nil {
		return err
	}

	nonce, err := corecontracts.GetAccountNonce(ch, agentID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	nonceResponse := &models.AccountNonceResponse{
		Nonce: iotago.EncodeUint64(nonce),
	}

	return e.JSON(http.StatusOK, nonceResponse)
}

func (c *Controller) getNFTData(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	nftID, err := params.DecodeNFTID(e)
	if err != nil {
		return err
	}

	nftData, err := corecontracts.GetNFTData(ch, *nftID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	nftDataResponse := isc.NFTToJSONObject(nftData)

	return e.JSON(http.StatusOK, nftDataResponse)
}

func (c *Controller) getNativeTokenIDRegistry(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	registries, err := corecontracts.GetNativeTokenIDRegistry(ch, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	nativeTokenIDRegistryResponse := &models.NativeTokenIDRegistryResponse{
		NativeTokenRegistryIDs: make([]string, len(registries)),
	}

	for k, v := range registries {
		nativeTokenIDRegistryResponse.NativeTokenRegistryIDs[k] = v.String()
	}

	return e.JSON(http.StatusOK, nativeTokenIDRegistryResponse)
}

func (c *Controller) getFoundryOutput(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	serialNumber, err := params.DecodeUInt(e, "serialNumber")
	if err != nil {
		return err
	}

	foundryOutput, err := corecontracts.GetFoundryOutput(ch, uint32(serialNumber), e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	foundryOutputID, err := foundryOutput.ID()
	if err != nil {
		return apierrors.InvalidPropertyError("FoundryOutput.ID", err)
	}

	foundryOutputResponse := &models.FoundryOutputResponse{
		FoundryID: foundryOutputID.ToHex(),
		Assets: models.AssetsResponse{
			BaseTokens:   iotago.EncodeUint64(foundryOutput.Amount),
			NativeTokens: isc.NativeTokensToJSONObject(foundryOutput.NativeTokens),
		},
	}

	return e.JSON(http.StatusOK, foundryOutputResponse)
}
