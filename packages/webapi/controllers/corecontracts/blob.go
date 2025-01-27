package corecontracts

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/controllers/controllerutils"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/corecontracts"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/params"
)

type Blob struct {
	Hash string `json:"hash" swagger:"required"`
	Size uint32 `json:"size" swagger:"required,min(1)"`
}

type BlobListResponse struct {
	Blobs []Blob
}

type BlobValueResponse struct {
	ValueData string `json:"valueData" swagger:"required,desc(The blob data (Hex))"`
}

func (c *Controller) getBlobValue(e echo.Context) error {
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	blobHash, err := params.DecodeBlobHash(e)
	if err != nil {
		return err
	}

	fieldKey := e.Param(params.ParamFieldKey)

	blobValueBytes, err := corecontracts.GetBlobValue(ch, *blobHash, fieldKey, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	blobValueResponse := &BlobValueResponse{
		ValueData: iotago.EncodeHex(blobValueBytes),
	}

	return e.JSON(http.StatusOK, blobValueResponse)
}

type BlobInfoResponse struct {
	Fields map[string]uint32 `json:"fields" swagger:"required,min(1)"`
}

func (c *Controller) getBlobInfo(e echo.Context) error {
	fmt.Println("GET BLOB INFO")

	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return c.handleViewCallError(err)
	}

	blobHash, err := params.DecodeBlobHash(e)
	if err != nil {
		return err
	}

	blobInfo, ok, err := corecontracts.GetBlobInfo(ch, *blobHash, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err)
	}

	fmt.Printf("GET BLOB INFO: ok:%v, err:%v", ok, err)

	blobInfoResponse := &BlobInfoResponse{
		Fields: map[string]uint32{},
	}

	// TODO: Validate this logic. Should an empty blob info result in a different http error code?
	if !ok {
		fmt.Printf("GET BLOB INFO return empty blobInfoResponse")
		return e.JSON(http.StatusOK, blobInfoResponse)
	}

	for k, v := range blobInfo {
		blobInfoResponse.Fields[k] = v
	}

	return e.JSON(http.StatusOK, blobInfoResponse)
}
