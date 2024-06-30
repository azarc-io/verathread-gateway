package internal

import (
	"github.com/azarc-io/verathread-gateway/internal/api"
	"github.com/azarc-io/verathread-gateway/internal/cache"
	federation2 "github.com/azarc-io/verathread-gateway/internal/federation"
	"net/http"
	"net/url"
	"strings"
	"time"

	error2 "github.com/azarc-io/verathread-gateway/internal/error"

	"github.com/azarc-io/verathread-gateway/internal/config"
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
	Domain struct {
		log          zerolog.Logger // default logger
		opts         *config.APIGatewayOptions
		http         httpuc.HttpUseCase
		httpClient   *http.Client
		sdlMap       map[string]string
		services     []*federation2.ServiceConfig
		federation   *federation2.Federation
		moduleMap    map[string]*ProxyTarget
		ready        bool
		is           api.InternalService
		httpInternal httpuc.HttpUseCase
		cache        *cache.ProjectCache
		registry     *RegistrationActor
	}
)

/************************************************************************/
/* LIFECYCLE
/************************************************************************/

func (d *Domain) Init() error {
	// handles caching apps to reduce db hits
	d.cache = cache.NewProjectCache(d.log)

	healthz.Register("gateway", time.Second*1, func() error {
		if !d.ready {
			return error2.ErrGatewayNotReady
		}
		return nil
	})

	// create service to handle inbound requests
	d.is = NewService(d.opts, d.log, d.cache)

	// register app registration actor
	if err := d.createAppRegistryHandler(); err != nil {
		return err
	}

	return nil
}

func (d *Domain) Stop() error {
	return nil
}

func (d *Domain) Start() error {
	// create http server
	d.http = httpuc.NewHttpUseCase(
		httpuc.WithHttpConfig(d.opts.Config.HTTP),
		httpuc.WithLogger(d.log),
	)

	// route all unhandled requests to echo
	mux := d.opts.DaprUseCase.Mux()
	mux.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		d.http.Server().ServeHTTP(writer, request)
	})

	// health check
	d.http.Server().GET("/health", echo.WrapHandler(healthz.Handler()))

	// register the shell app route
	d.registerShellAppRoute()

	// graphql federation
	if err := d.registerGraphqlRoute(); err != nil {
		return err
	}

	// resource proxy
	d.registerProxyRouter()

	// flag service is ready so health starts reporting ok status
	d.ready = true

	return nil
}

/************************************************************************/
/* SHELL APP
/************************************************************************/

// registerShellAppRoute serves up the shell app
func (d *Domain) registerShellAppRoute() {
	e := d.http.Server()

	if d.opts.Config.WebProxy != "" {
		d.log.Info().Msgf("serving files from: %s", d.opts.Config.WebProxy)

		_url, err := url.Parse(d.opts.Config.WebProxy)
		if err != nil {
			panic(err)
		}

		tgt := &ProxyTarget{
			Name:         "shell",
			URL:          _url,
			Meta:         nil,
			RegexRewrite: nil,
		}

		grp := d.http.Server().Group("")
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
	} else if d.opts.Config.WebDir != "" {
		d.log.Info().Msgf("serving files from: %s", d.opts.Config.WebDir)

		g1 := e.Group("")
		g1.Use(
			middleware.GzipWithConfig(middleware.GzipConfig{
				Skipper: func(c echo.Context) bool {
					ct := c.Response().Header().Get(echo.HeaderContentType)
					return ct != "text/css" && ct != "application/javascript"
				},
			}),
			middleware.StaticWithConfig(middleware.StaticConfig{
				Root:  d.opts.Config.WebDir,
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

func (d *Domain) registerGraphqlRoute() error {
	// transform any statically registered routes provided through configuration
	for name, service := range d.opts.Config.Services {
		d.services = append(d.services, &federation2.ServiceConfig{
			Name:     name,
			URL:      service.Gql,
			WS:       service.GqlWs,
			Fallback: nil,
		})
	}

	// register the application gateways own graphql endpoint with the federation server,
	// so we can serve up information about apps, navigation etc.
	if err := d.registerGqlApi(); err != nil {
		return err
	}

	// create the federation gateway
	d.federation = federation2.New(d.http, d.log, d.services, d.opts.WardenUseCase)

	return nil
}

/************************************************************************/
/* WEB APP PROXIES
/************************************************************************/

func (d *Domain) registerProxyRouter() {
	grp := d.http.Server().Group("/module/:name")
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

			if tgt, ok := d.moduleMap[module]; ok {
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

func NewGateway(opts ...config.APIGatewayOption) *Domain {
	g := &Domain{
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
