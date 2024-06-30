package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	federation2 "github.com/azarc-io/verathread-gateway/internal/federation"
	"github.com/azarc-io/verathread-gateway/internal/gql/graph"
	"github.com/azarc-io/verathread-gateway/internal/gql/graph/resolvers"
	httpuc "github.com/azarc-io/verathread-next-common/usecase/http"
	netutil "github.com/azarc-io/verathread-next-common/util/net"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

// registerGqlApi registers graphql api handler
func (d *Domain) registerGqlApi() error {
	// get a free port
	port, err := netutil.GetFreeTCPPort()
	if err != nil {
		return fmt.Errorf("could not find a free port for internal server: %w", err)
	}

	// create internal http server so we can federate it
	// TODO would be good to find a way to do this without having to create an internal server
	d.httpInternal = httpuc.NewHttpUseCase(
		httpuc.WithHttpConfig(&httpuc.ConfigBindHttp{
			Address: "127.0.0.1",
			Port:    port,
		}),
		httpuc.WithLogger(d.log),
	)

	resolver := &resolvers.Resolver{
		Opts:            d.opts,
		InternalService: d.is,
	}

	c := graph.Config{
		Resolvers: resolver,
	}

	sch := graph.NewExecutableSchema(c)
	srv := handler.New(sch)
	srv.Use(extension.FixedComplexityLimit(70))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: time.Second * 10,
		PingPongInterval:      time.Second * 20,
		InitTimeout:           time.Second * 20,
		Upgrader: websocket.Upgrader{
			HandshakeTimeout: time.Second * 10,
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				d.log.Error().Err(reason).Msgf("error during ws upgrade")
			},
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		//InitFunc: gqlMw.WebSocketInit,
		ErrorFunc: func(ctx context.Context, err error) {
			if !websocket.IsCloseError(err) &&
				!errors.Is(err, websocket.ErrCloseSent) &&
				!strings.Contains(err.Error(), "going away") &&
				!strings.Contains(err.Error(), "close sent") &&
				!strings.Contains(err.Error(), "websocket connection closed") {
				d.log.Error().Err(err).Msgf("error during ws connection")
			}
		},
	})

	d.httpInternal.Server().Any("/internal/graphql", echo.WrapHandler(srv))

	// register this service as a static route
	uri := fmt.Sprintf("http://127.0.0.1:%d/internal/graphql", port)
	d.services = append(d.services, &federation2.ServiceConfig{
		Name:     "apps",
		URL:      uri,
		WS:       uri,
		Fallback: nil,
	})

	return d.httpInternal.Start()
}

// createAppRegistryHandler registers the actor that other cluster members can use to register themselves
func (d *Domain) createAppRegistryHandler() error {
	d.registry = newAppRegistrationHandler(d.cache, d.opts, d.log, d.is)
	return d.registry.RegisterHandlers()
}
