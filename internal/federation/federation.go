package federation

import (
	"context"
	http2 "github.com/azarc-io/verathread-gateway/internal/federation/http"
	"net/http"
	"time"

	httpuc "github.com/azarc-io/verathread-next-common/usecase/http"
	wardenuc "github.com/azarc-io/verathread-next-common/usecase/warden"
	"github.com/gobwas/ws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/playground"
)

var (
	maxIdleConnections = 10
	idleConnTimeout    = 30 * time.Second
	pollingInterval    = 30 * time.Second
)

type Federation struct {
	gateway           *Gateway
	datasourceWatcher *DatasourcePollerPoller
}

func (f Federation) RegisterAPI(pkg string, url, wsURL string) {
	f.datasourceWatcher.addService(pkg, url, wsURL)
}

func (f Federation) DeRegisterAPI(pkg string) {
	f.datasourceWatcher.removeService(pkg)
}

func New(huc httpuc.HttpUseCase, log zerolog.Logger, services []*ServiceConfig, wuc wardenuc.ClusterWardenUseCase) *Federation {
	graphqlEndpoint := "/graphql"
	playgroundURLPrefix := "/playground"

	tr := &http.Transport{
		MaxIdleConns:       maxIdleConnections,
		IdleConnTimeout:    idleConnTimeout,
		DisableCompression: false,
	}
	httpClient := &http.Client{Transport: wuc.NewOtelTransport(tr)}

	upgrader := &ws.DefaultHTTPUpgrader
	upgrader.Header = http.Header{}
	upgrader.Header.Add("Sec-Websocket-Protocol", "graphql-ws")

	datasourceWatcher := NewDatasourcePoller(httpClient, log, DatasourcePollerConfig{
		Services:        services,
		PollingInterval: pollingInterval,
	})

	p := playground.New(playground.Config{
		PathPrefix:                      "",
		PlaygroundPath:                  playgroundURLPrefix,
		GraphqlEndpointPath:             graphqlEndpoint,
		GraphQLSubscriptionEndpointPath: graphqlEndpoint,
	})

	handlers, err := p.Handlers()
	if err != nil {
		log.Fatal().Err(err).Msgf("configure handlers")
		return nil
	}

	for i := range handlers {
		huc.Server().Any(handlers[i].Path, echo.WrapHandler(handlers[i].Handler))
	}

	enableART := true

	var gqlHandlerFactory HandlerFactoryFn = func(schema *graphql.Schema, engine *engine.ExecutionEngine) http.Handler {
		return http2.NewGraphqlHTTPHandler(schema, engine, upgrader, log, enableART)
	}

	gateway := NewGqlGateway(context.Background(), gqlHandlerFactory, httpClient, log)
	datasourceWatcher.Register(gateway)

	if len(services) > 0 {
		go datasourceWatcher.Run(context.Background())
		gateway.Ready()
	}

	huc.Server().Any("/graphql", func(c echo.Context) error {
		w := c.Response()
		r := c.Request()

		// handover to the gateway
		gateway.ServeHTTP(w, r)

		return nil
	}, wuc.HttpMiddleware)

	return &Federation{
		gateway:           gateway,
		datasourceWatcher: datasourceWatcher,
	}
}
