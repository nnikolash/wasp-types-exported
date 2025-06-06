package webapi

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
	"go.uber.org/dig"
	websocketserver "nhooyr.io/websocket"

	"github.com/iotaledger/hive.go/app"
	"github.com/iotaledger/hive.go/app/configuration"
	"github.com/iotaledger/hive.go/app/shutdown"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/web/websockethub"
	"github.com/iotaledger/inx-app/pkg/httpserver"
	"github.com/nnikolash/wasp-types-exported/packages/authentication"
	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/chains"
	"github.com/nnikolash/wasp-types-exported/packages/daemon"
	"github.com/nnikolash/wasp-types-exported/packages/dkg"
	"github.com/nnikolash/wasp-types-exported/packages/evm/jsonrpc"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/metrics"
	"github.com/nnikolash/wasp-types-exported/packages/peering"
	"github.com/nnikolash/wasp-types-exported/packages/publisher"
	"github.com/nnikolash/wasp-types-exported/packages/registry"
	"github.com/nnikolash/wasp-types-exported/packages/users"
	"github.com/nnikolash/wasp-types-exported/packages/webapi"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/apierrors"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/controllers/controllerutils"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/websocket"
)

func init() {
	Component = &app.Component{
		Name:             "WebAPI",
		DepsFunc:         func(cDeps dependencies) { deps = cDeps },
		Params:           params,
		InitConfigParams: initConfigParams,
		IsEnabled:        func(_ *dig.Container) bool { return ParamsWebAPI.Enabled },
		Provide:          provide,
		Run:              run,
	}
}

var (
	Component *app.Component
	deps      dependencies
)

const (
	broadcastQueueSize            = 20000
	clientSendChannelSize         = 1000
	maxWebsocketMessageSize int64 = 510
)

type dependencies struct {
	dig.In

	EchoSwagger        echoswagger.ApiRoot `name:"webapiServer"`
	WebsocketHub       *websockethub.Hub   `name:"websocketHub"`
	NodeConnection     chain.NodeConnection
	WebsocketPublisher *websocket.Service `name:"websocketService"`
}

func initConfigParams(c *dig.Container) error {
	type cfgResult struct {
		dig.Out
		WebAPIBindAddress string `name:"webAPIBindAddress"`
	}

	if err := c.Provide(func() cfgResult {
		return cfgResult{
			WebAPIBindAddress: ParamsWebAPI.BindAddress,
		}
	}); err != nil {
		Component.LogPanic(err)
	}

	return nil
}

//nolint:funlen
func NewEcho(params *ParametersWebAPI, metrics *metrics.ChainMetricsProvider, log *logger.Logger) *echo.Echo {
	e := httpserver.NewEcho(
		log,
		nil,
		ParamsWebAPI.DebugRequestLoggerEnabled,
	)

	e.Server.ReadTimeout = params.Limits.ReadTimeout
	e.Server.WriteTimeout = params.Limits.WriteTimeout

	e.HidePort = true
	e.HTTPErrorHandler = apierrors.HTTPErrorHandler()

	webapi.ConfirmedStateLagThreshold = params.Limits.ConfirmedStateLagThreshold
	authentication.DefaultJWTDuration = params.Auth.JWTConfig.Duration

	e.Pre(middleware.RemoveTrailingSlash())

	// publish metrics to prometheus component (that exposes a separate http server on another port)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Path(), "/chains/") {
				// ignore metrics for all requests not related to "chains/<chainID>""
				return next(c)
			}
			start := time.Now()
			err := next(c)

			status := c.Response().Status
			if err != nil {
				var httpError *echo.HTTPError
				if errors.As(err, &httpError) {
					status = httpError.Code
				}
				if status == 0 || status == http.StatusOK {
					status = http.StatusInternalServerError
				}
			}

			chainID, ok := c.Get(controllerutils.EchoContextKeyChainID).(isc.ChainID)
			if !ok {
				return err
			}

			operation, ok := c.Get(controllerutils.EchoContextKeyOperation).(string)
			if !ok {
				return err
			}
			metrics.GetChainMetrics(chainID).WebAPI.WebAPIRequest(operation, status, time.Since(start))
			return err
		}
	})

	// timeout middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeoutCtx, cancel := context.WithTimeout(c.Request().Context(), params.Limits.Timeout)
			defer cancel()

			c.SetRequest(c.Request().WithContext(timeoutCtx))

			return next(c)
		}
	})

	// Middleware to unescape any supplied path (/path/foo%40bar/) parameter
	// Query parameters (?name=foo%40bar) get unescaped by default.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			escapedPathParams := c.ParamValues()
			unescapedPathParams := make([]string, len(escapedPathParams))

			for i, param := range escapedPathParams {
				unescapedParam, err := url.PathUnescape(param)

				if err != nil {
					unescapedPathParams[i] = param
				} else {
					unescapedPathParams[i] = unescapedParam
				}
			}

			c.SetParamValues(unescapedPathParams...)

			return next(c)
		}
	})

	e.Use(middleware.BodyLimit(params.Limits.MaxBodyLength))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${remote_ip} ${method} ${uri} ${status} error="${error}"` + "\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowCredentials: true,
	}))

	return e
}

func CreateEchoSwagger(e *echo.Echo, version string) echoswagger.ApiRoot {
	echoSwagger := echoswagger.New(e, "/doc", &echoswagger.Info{
		Title:       "Wasp API",
		Description: "REST API for the Wasp node",
		Version:     version,
	})

	echoSwagger.AddSecurityAPIKey("Authorization", "JWT Token", echoswagger.SecurityInHeader).
		SetExternalDocs("Find out more about Wasp", "https://wiki.iota.org/smart-contracts/overview").
		SetUI(echoswagger.UISetting{DetachSpec: false, HideTop: false}).
		SetScheme("http", "https")

	echoSwagger.SetRequestContentType(echo.MIMEApplicationJSON)
	echoSwagger.SetResponseContentType(echo.MIMEApplicationJSON)

	return echoSwagger
}

//nolint:funlen
func provide(c *dig.Container) error {
	type webapiServerDeps struct {
		dig.In

		AppInfo                     *app.Info
		AppConfig                   *configuration.Configuration `name:"appConfig"`
		ShutdownHandler             *shutdown.ShutdownHandler
		APICacheTTL                 time.Duration `name:"apiCacheTTL"`
		Chains                      *chains.Chains
		ChainMetricsProvider        *metrics.ChainMetricsProvider
		ChainRecordRegistryProvider registry.ChainRecordRegistryProvider
		DKShareRegistryProvider     registry.DKShareRegistryProvider
		NodeIdentityProvider        registry.NodeIdentityProvider
		NetworkProvider             peering.NetworkProvider       `name:"networkProvider"`
		TrustedNetworkManager       peering.TrustedNetworkManager `name:"trustedNetworkManager"`
		Node                        *dkg.Node
		UserManager                 *users.UserManager
		Publisher                   *publisher.Publisher
	}

	type webapiServerResult struct {
		dig.Out

		Echo               *echo.Echo          `name:"webapiEcho"`
		EchoSwagger        echoswagger.ApiRoot `name:"webapiServer"`
		WebsocketHub       *websockethub.Hub   `name:"websocketHub"`
		WebsocketPublisher *websocket.Service  `name:"websocketService"`
	}

	if err := c.Provide(func(deps webapiServerDeps) webapiServerResult {
		e := NewEcho(ParamsWebAPI, deps.ChainMetricsProvider, Component.Logger())

		echoSwagger := CreateEchoSwagger(e, deps.AppInfo.Version)
		websocketOptions := websocketserver.AcceptOptions{
			InsecureSkipVerify: true,
			// Disable compression due to incompatibilities with the latest Safari browsers:
			// https://github.com/tilt-dev/tilt/issues/4746
			CompressionMode: websocketserver.CompressionDisabled,
		}

		logger := Component.App().NewLogger("WebAPI/v2")

		hub := websockethub.NewHub(Component.Logger(), &websocketOptions, broadcastQueueSize, clientSendChannelSize, maxWebsocketMessageSize)

		websocketService := websocket.NewWebsocketService(logger, hub, []publisher.ISCEventType{
			publisher.ISCEventKindNewBlock,
			publisher.ISCEventKindReceipt,
			publisher.ISCEventIssuerVM,
			publisher.ISCEventKindBlockEvents,
		}, deps.Publisher, websocket.WithMaxTopicSubscriptionsPerClient(ParamsWebAPI.Limits.MaxTopicSubscriptionsPerClient))

		if ParamsWebAPI.DebugRequestLoggerEnabled {
			echoSwagger.Echo().Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
				logger.Debugf("API Dump: Request=%q, Response=%q", reqBody, resBody)
			}))
		}

		webapi.Init(
			logger,
			echoSwagger,
			deps.AppInfo.Version,
			deps.AppConfig,
			deps.NetworkProvider,
			deps.TrustedNetworkManager,
			deps.UserManager,
			deps.ChainRecordRegistryProvider,
			deps.DKShareRegistryProvider,
			deps.NodeIdentityProvider,
			func() *chains.Chains {
				return deps.Chains
			},
			func() *dkg.Node {
				return deps.Node
			},
			deps.ShutdownHandler,
			deps.ChainMetricsProvider,
			ParamsWebAPI.Auth,
			deps.APICacheTTL,
			websocketService,
			ParamsWebAPI.IndexDbPath,
			ParamsWebAPI.AccountDumpsPath,
			deps.Publisher,
			jsonrpc.NewParameters(
				ParamsWebAPI.Limits.Jsonrpc.MaxBlocksInLogsFilterRange,
				ParamsWebAPI.Limits.Jsonrpc.MaxLogsInResult,
				ParamsWebAPI.Limits.Jsonrpc.WebsocketRateLimitMessagesPerSecond,
				ParamsWebAPI.Limits.Jsonrpc.WebsocketRateLimitBurst,
				ParamsWebAPI.Limits.Jsonrpc.WebsocketConnectionCleanupDuration,
				ParamsWebAPI.Limits.Jsonrpc.WebsocketClientBlockDuration,
			),
		)

		return webapiServerResult{
			EchoSwagger:        echoSwagger,
			WebsocketHub:       hub,
			WebsocketPublisher: websocketService,
		}
	}); err != nil {
		Component.LogPanic(err)
	}

	return nil
}

func run() error {
	Component.LogInfof("Starting %s server ...", Component.Name)
	if err := Component.Daemon().BackgroundWorker(Component.Name, func(ctx context.Context) {
		Component.LogInfof("Starting %s server ...", Component.Name)
		if err := deps.NodeConnection.WaitUntilInitiallySynced(ctx); err != nil {
			Component.LogErrorf("failed to start %s, waiting for L1 node to become sync failed, error: %s", err.Error())
			return
		}

		Component.LogInfof("Starting %s server ... done", Component.Name)

		go func() {
			deps.EchoSwagger.Echo().Server.BaseContext = func(_ net.Listener) context.Context {
				// set BaseContext to be the same as the plugin, so that requests being processed don't hang the shutdown procedure
				return ctx
			}

			Component.LogInfof("You can now access the WebAPI using: http://%s", ParamsWebAPI.BindAddress)
			if err := deps.EchoSwagger.Echo().Start(ParamsWebAPI.BindAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
				Component.LogWarnf("Stopped %s server due to an error (%s)", Component.Name, err)
			}
		}()

		<-ctx.Done()

		Component.LogInfof("Stopping %s server ...", Component.Name)

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCtxCancel()

		//nolint:contextcheck // false positive
		if err := deps.EchoSwagger.Echo().Shutdown(shutdownCtx); err != nil {
			Component.LogWarn(err)
		}

		Component.LogInfof("Stopping %s server ... done", Component.Name)
	}, daemon.PriorityWebAPI); err != nil {
		Component.LogPanicf("failed to start worker: %s", err)
	}

	if err := Component.Daemon().BackgroundWorker("WebAPI[WS]", func(ctx context.Context) {
		unhook := deps.WebsocketPublisher.EventHandler().AttachToEvents()
		defer unhook()

		deps.WebsocketHub.Run(ctx)
		Component.LogInfo("Stopping WebAPI[WS]")
	}, daemon.PriorityWebAPI); err != nil {
		Component.LogPanicf("failed to start worker: %s", err)
	}

	return nil
}
