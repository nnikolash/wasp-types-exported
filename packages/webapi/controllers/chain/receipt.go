package chain

import (
	"github.com/labstack/echo/v4"

	"github.com/nnikolash/wasp-types-exported/packages/webapi/controllers/controllerutils"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/controllers/corecontracts"
)

func (c *Controller) getReceipt(e echo.Context) error {
	controllerutils.SetOperation(e, "get_receipt")
	return corecontracts.GetRequestReceipt(e, c.chainService)
}
