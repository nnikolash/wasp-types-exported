package webapi_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/nnikolash/wasp-types-exported/components/webapi"
	"github.com/nnikolash/wasp-types-exported/packages/authentication"
)

func TestInternalServerErrors(t *testing.T) {
	// start a webserver with a test log
	logCore, logObserver := observer.New(zapcore.DebugLevel)
	log := zap.New(logCore)

	e := webapi.NewEcho(&webapi.ParametersWebAPI{
		Enabled:     true,
		BindAddress: ":9999",
		Auth:        authentication.AuthConfiguration{},
		Limits: webapi.ParametersWebAPILimits{
			Timeout:                        time.Minute,
			ReadTimeout:                    time.Minute,
			WriteTimeout:                   time.Minute,
			MaxBodyLength:                  "1M",
			MaxTopicSubscriptionsPerClient: 0,
			ConfirmedStateLagThreshold:     2,
			Jsonrpc:                        webapi.ParametersJSONRPC{},
		},
		DebugRequestLoggerEnabled: true,
	},
		nil,
		log.Sugar(),
	)

	// Add an endpoint that just panics with "foobar" and start the server
	exceptionText := "foobar"
	e.GET("/test", func(c echo.Context) error { panic(exceptionText) })
	go func() {
		err := e.Start(":9999")
		require.ErrorIs(t, http.ErrServerClosed, err)
	}()
	defer e.Shutdown(context.Background())

	// query the endpoint
	req, err := http.NewRequest(http.MethodGet, "http://localhost:9999/test", http.NoBody)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	res.Body.Close()

	// assert the exception is not present in the response (prevent leaking errors)
	require.Equal(t, res.StatusCode, http.StatusInternalServerError)
	require.NotContains(t, string(resBody), exceptionText)

	// assert the exception is logged
	logEntries := logObserver.All()
	require.Len(t, logEntries, 1)
	require.Contains(t, logEntries[0].Message, exceptionText)
}
