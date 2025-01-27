package chain

import (
	"net/http"

	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/apierrors"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/common"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/controllers/controllerutils"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/models"

	"github.com/labstack/echo/v4"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
)

func (c *Controller) executeCallView(e echo.Context) error {
	controllerutils.SetOperation(e, "call_view")
	ch, _, err := controllerutils.ChainFromParams(e, c.chainService)
	if err != nil {
		return err
	}

	var callViewRequest models.ContractCallViewRequest
	if err = e.Bind(&callViewRequest); err != nil {
		return apierrors.InvalidPropertyError("body", err)
	}

	// Get contract and function. The request model supports HName and common string names. HNames are preferred.
	var contractHName isc.Hname
	var functionHName isc.Hname

	if callViewRequest.ContractHName == "" {
		contractHName = isc.Hn(callViewRequest.ContractName)
	} else {
		contractHName, err = isc.HnameFromString(callViewRequest.ContractHName)
		if err != nil {
			return apierrors.InvalidPropertyError("contractHName", err)
		}
	}

	if callViewRequest.FunctionHName == "" {
		functionHName = isc.Hn(callViewRequest.FunctionName)
	} else {
		functionHName, err = isc.HnameFromString(callViewRequest.FunctionHName)
		if err != nil {
			return apierrors.InvalidPropertyError("contractHName", err)
		}
	}

	args, err := dict.FromJSONDict(callViewRequest.Arguments)
	if err != nil {
		return apierrors.InvalidPropertyError("arguments", err)
	}

	result, err := common.CallView(ch, contractHName, functionHName, args, callViewRequest.Block)
	if err != nil {
		return apierrors.ContractExecutionError(err)
	}

	return e.JSON(http.StatusOK, result)
}
