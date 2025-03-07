package metrics

import (
	"net/http"

	"github.com/pangpanglabs/echoswagger/v2"

	"github.com/nnikolash/wasp-types-exported/packages/authentication"
	"github.com/nnikolash/wasp-types-exported/packages/authentication/shared/permissions"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/interfaces"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/models"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/params"
)

type Controller struct {
	chainService   interfaces.ChainService
	metricsService interfaces.MetricsService
}

func NewMetricsController(chainService interfaces.ChainService, metricsService interfaces.MetricsService) interfaces.APIController {
	return &Controller{
		chainService:   chainService,
		metricsService: metricsService,
	}
}

func (c *Controller) Name() string {
	return "metrics"
}

func (c *Controller) RegisterPublic(publicAPI echoswagger.ApiGroup, mocker interfaces.Mocker) {
}

func (c *Controller) RegisterAdmin(adminAPI echoswagger.ApiGroup, mocker interfaces.Mocker) {
	adminAPI.GET("metrics/node/messages", c.getNodeMessageMetrics, authentication.ValidatePermissions([]string{permissions.Read})).
		AddResponse(http.StatusOK, "A list of all available metrics.", models.NodeMessageMetrics{}, nil).
		SetOperationId("getNodeMessageMetrics").
		SetSummary("Get accumulated message metrics.")

	adminAPI.GET("metrics/chain/:chainID/messages", c.getChainMessageMetrics, authentication.ValidatePermissions([]string{permissions.Read})).
		AddParamPath("", params.ParamChainID, params.DescriptionChainID).
		AddResponse(http.StatusNotFound, "Chain not found", nil, nil).
		AddResponse(http.StatusOK, "A list of all available metrics.", models.ChainMessageMetrics{}, nil).
		SetOperationId("getChainMessageMetrics").
		SetSummary("Get chain specific message metrics.")

	adminAPI.GET("metrics/chain/:chainID/workflow", c.getChainWorkflowMetrics, authentication.ValidatePermissions([]string{permissions.Read})).
		AddParamPath("", params.ParamChainID, params.DescriptionChainID).
		AddResponse(http.StatusNotFound, "Chain not found", nil, nil).
		AddResponse(http.StatusOK, "A list of all available metrics.", mocker.Get(models.ConsensusWorkflowMetrics{}), nil).
		SetOperationId("getChainWorkflowMetrics").
		SetSummary("Get chain workflow metrics.")

	adminAPI.GET("metrics/chain/:chainID/pipe", c.getChainPipeMetrics, authentication.ValidatePermissions([]string{permissions.Read})).
		AddParamPath("", params.ParamChainID, params.DescriptionChainID).
		AddResponse(http.StatusNotFound, "Chain not found", nil, nil).
		AddResponse(http.StatusOK, "A list of all available metrics.", mocker.Get(models.ConsensusPipeMetrics{}), nil).
		SetOperationId("getChainPipeMetrics").
		SetSummary("Get chain pipe event metrics.")
}
