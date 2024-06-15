package federation

import (
	"context"
	"net/http"
	"sync"

	"github.com/azarc-io/verathread-gateway/internal/gateway/federation/logger"
	"github.com/rs/zerolog"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

var maxConcurrency = 1024

type (
	Gateway struct {
		engineCtx         context.Context
		gqlHandlerFactory HandlerFactory
		httpClient        *http.Client
		log               zerolog.Logger
		mu                *sync.Mutex
		readyCh           chan struct{}
		readyOnce         *sync.Once
		gqlHandler        http.Handler
	}

	DataSourceObserver interface {
		UpdateDataSources(subgraphsConfigs []engine.SubgraphConfiguration)
	}

	DataSourceSubject interface {
		Register(observer DataSourceObserver)
	}

	HandlerFactory interface {
		Make(schema *graphql.Schema, engine *engine.ExecutionEngine) http.Handler
	}

	HandlerFactoryFn func(schema *graphql.Schema, engine *engine.ExecutionEngine) http.Handler
)

func (h HandlerFactoryFn) Make(schema *graphql.Schema, engine *engine.ExecutionEngine) http.Handler {
	return h(schema, engine)
}

func (g *Gateway) UpdateDataSources(subgraphsConfigs []engine.SubgraphConfiguration) {
	engineConfigFactory := engine.NewFederationEngineConfigFactory(
		g.engineCtx,
		subgraphsConfigs,
		engine.WithFederationHttpClient(g.httpClient),
	)

	engineConfig, err := engineConfigFactory.BuildEngineConfiguration()
	if err != nil {
		g.log.Error().Err(err).Msgf("get engine config")
		return
	}

	executionEngine, err := engine.NewExecutionEngine(
		g.engineCtx,
		logger.ZapLogger(g.log),
		engineConfig,
		resolve.ResolverOptions{
			MaxConcurrency:               maxConcurrency,
			PropagateSubgraphErrors:      true,
			PropagateSubgraphStatusCodes: false,
			SubgraphErrorPropagationMode: resolve.SubgraphErrorPropagationModePassThrough,
		})
	if err != nil {
		g.log.Error().Err(err).Msgf("create engine")
		return
	}

	g.mu.Lock()
	g.gqlHandler = g.gqlHandlerFactory.Make(engineConfig.Schema(), executionEngine)
	g.mu.Unlock()

	g.readyOnce.Do(func() { close(g.readyCh) })
}

func (g *Gateway) Ready() {
	<-g.readyCh
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	handler := g.gqlHandler
	g.mu.Unlock()

	if handler != nil {
		handler.ServeHTTP(w, r)
	}
}

func NewGqlGateway(ctx context.Context, factory HandlerFactory, client *http.Client, log zerolog.Logger) *Gateway {
	return &Gateway{
		engineCtx:         ctx,
		gqlHandlerFactory: factory,
		httpClient:        client,
		log:               log,

		mu:        &sync.Mutex{},
		readyCh:   make(chan struct{}),
		readyOnce: &sync.Once{},
	}
}
