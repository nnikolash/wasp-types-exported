package webapi

import (
	_ "embed"
	"fmt"

	"github.com/pangpanglabs/echoswagger/v2"

	"github.com/nnikolash/wasp-types-exported/packages/webapi/websocket"
)

func addWebSocketEndpoint(e echoswagger.ApiRoot, websocketPublisher *websocket.Service) {
	e.GET(fmt.Sprintf("/v%d/ws", APIVersion), websocketPublisher.ServeHTTP).
		SetSummary("The websocket connection service")
}
