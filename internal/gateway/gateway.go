package gateway

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	error2 "github.com/azarc-io/verathread-gateway/internal/error"

	"github.com/azarc-io/verathread-gateway/internal/config"
	"github.com/azarc-io/verathread-gateway/internal/gateway/federation"
	httpuc "github.com/azarc-io/verathread-next-common/usecase/http"
	"github.com/azarc-io/verathread-next-common/util/healthz"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/************************************************************************/
/* TYPES
/************************************************************************/

type (
	Gateway struct {
		log        zerolog.Logger // default logger
		opts       *config.APIGatewayOptions
		http       httpuc.HttpUseCase
		httpClient *http.Client
		sdlMap     map[string]string
		services   []*federation.ServiceConfig
		federation *federation.Federation
		moduleMap  map[string]*ProxyTarget
		ready      bool
	}

	AppAddedOrUpdatedEvent struct {
		APIEndpoint   string       `json:"apiUrl"`
		APIWsEndpoint string       `json:"apiWsUrl"`
		Package       string       `json:"package"`
		ProxyAPI      bool         `json:"proxyApi"`
		Name          string       `json:"name"`
		Version       string       `json:"version"`
		Available     bool         `json:"available"`
		Navigation    []*AppModule `json:"navigation"`
	}

	AppRemovedEvent struct {
		Package    string       `json:"package"`
		Name       string       `json:"name"`
		Version    string       `json:"version"`
		Navigation []*AppModule `json:"navigation"`
	}

	AppModule struct {
		ID                      string            `json:"id"`
		Proxy                   bool              `json:"proxy"`
		BaseURL                 string            `json:"baseUrl"`
		RemoteEntryRewriteRegEx map[string]string `json:"remoteEntryRewriteRegEx"`
		Slug                    string            `json:"slug"`
	}
)

/************************************************************************/
/* LIFECYCLE
/************************************************************************/

func (g *Gateway) Stop() error {
	return nil
}

func (g *Gateway) Init() error {
	healthz.Register("gateway", time.Second*1, func() error {
		if !g.ready {
			return error2.ErrGatewayNotReady
		}
		return nil
	})

	return nil
}

func (g *Gateway) Start() error {
	// create http server
	g.http = httpuc.NewHttpUseCase(
		httpuc.WithHttpConfig(g.opts.Config.HTTP),
		httpuc.WithLogger(g.log),
	)

	// create dapr server service
	mux := g.opts.DaprUseCase.Mux()
	mux.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		g.http.Server().ServeHTTP(writer, request)
	})

	// health check
	g.http.Server().GET("/health", echo.WrapHandler(healthz.Handler()))

	// register the shell app route
	g.registerShellAppRoute()

	// graphql federation
	g.registerGraphqlRoute()

	// resource proxy
	g.registerProxyRouter()

	g.ready = true

	return nil
}

/************************************************************************/
/* SHELL APP
/************************************************************************/

// registerShellAppRoute serves up the shell app
func (g *Gateway) registerShellAppRoute() {
	e := g.http.Server()

	if g.opts.Config.WebProxy != "" {
		g.log.Info().Msgf("serving files from: %s", g.opts.Config.WebProxy)

		url, err := url.Parse(g.opts.Config.WebProxy)
		if err != nil {
			panic(err)
		}

		tgt := &ProxyTarget{
			Name:         "shell",
			URL:          url,
			Meta:         nil,
			RegexRewrite: nil,
		}

		grp := g.http.Server().Group("")
		grp.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				res := c.Response()

				if req.Header.Get(echo.HeaderXRealIP) == "" || c.Echo().IPExtractor != nil {
					req.Header.Set(echo.HeaderXRealIP, c.RealIP())
				}

				if req.Header.Get(echo.HeaderXForwardedProto) == "" {
					req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
				}

				// For HTTP, it is automatically set by Go HTTP reverse proxy.
				if c.IsWebSocket() && req.Header.Get(echo.HeaderXForwardedFor) == "" {
					req.Header.Set(echo.HeaderXForwardedFor, c.RealIP())
				}

				if err := rewriteURL(tgt.RegexRewrite, req); err != nil {
					return err
				}

				// Proxy
				switch {
				case c.IsWebSocket():
					proxyRaw(tgt, c).ServeHTTP(res, req)
				case req.Header.Get(echo.HeaderAccept) == "text/event-stream":
				default:
					log.Info().Msgf("proxy to %s%s", tgt.URL, req.URL)
					proxyHTTP(tgt, c).ServeHTTP(res, req)
				}

				return nil
			}
		})
	} else if g.opts.Config.WebDir != "" {
		g.log.Info().Msgf("serving files from: %s", g.opts.Config.WebDir)

		g1 := e.Group("")
		g1.Use(
			middleware.GzipWithConfig(middleware.GzipConfig{
				Skipper: func(c echo.Context) bool {
					ct := c.Response().Header().Get(echo.HeaderContentType)
					return ct != "text/css" && ct != "application/javascript"
				},
			}),
			middleware.StaticWithConfig(middleware.StaticConfig{
				Root:  g.opts.Config.WebDir,
				Index: "index.html",
				HTML5: true,
				Skipper: func(e echo.Context) bool {
					return strings.HasPrefix(e.Path(), "/tmp") ||
						strings.HasPrefix(e.Path(), "/api") ||
						strings.HasPrefix(e.Path(), "/graphql") ||
						strings.HasPrefix(e.Path(), "/query") ||
						strings.HasPrefix(e.Path(), "/health")
				},
			}),
		)
	}
}

/************************************************************************/
/* GRAPHQL ROUTING
/************************************************************************/

func (g *Gateway) registerGraphqlRoute() {
	for name, service := range g.opts.Config.Services {
		g.services = append(g.services, &federation.ServiceConfig{
			Name:     name,
			URL:      service.Gql,
			WS:       service.GqlWs,
			Fallback: nil,
		})
	}

	g.federation = federation.New(g.http, g.log, g.services, g.opts.WardenUseCase)
}

/************************************************************************/
/* WEB APP PROXIES
/************************************************************************/

func (g *Gateway) registerProxyRouter() {
	grp := g.http.Server().Group("/module/:name")
	grp.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			module := c.Param("name")

			if req.Header.Get(echo.HeaderXRealIP) == "" || c.Echo().IPExtractor != nil {
				req.Header.Set(echo.HeaderXRealIP, c.RealIP())
			}

			if req.Header.Get(echo.HeaderXForwardedProto) == "" {
				req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
			}

			// For HTTP, it is automatically set by Go HTTP reverse proxy.
			if c.IsWebSocket() && req.Header.Get(echo.HeaderXForwardedFor) == "" {
				req.Header.Set(echo.HeaderXForwardedFor, c.RealIP())
			}

			if tgt, ok := g.moduleMap[module]; ok {
				if err := rewriteURL(tgt.RegexRewrite, req); err != nil {
					return err
				}

				// Proxy
				switch {
				case c.IsWebSocket():
					proxyRaw(tgt, c).ServeHTTP(res, req)
				case req.Header.Get(echo.HeaderAccept) == "text/event-stream":
				default:
					log.Info().Msgf("proxy to %s%s", tgt.URL, req.URL)
					proxyHTTP(tgt, c).ServeHTTP(res, req)
				}
			}

			return nil
		}
	})
}

/************************************************************************/
/* FACTORY
/************************************************************************/

func NewGateway(opts ...config.APIGatewayOption) *Gateway {
	g := &Gateway{
		log:        log.With().Str("app", "gateway").Logger(),
		opts:       &config.APIGatewayOptions{},
		httpClient: http.DefaultClient,
		sdlMap:     make(map[string]string),
		moduleMap:  make(map[string]*ProxyTarget),
	}

	for _, opt := range opts {
		opt(g.opts)
	}

	return g
}
