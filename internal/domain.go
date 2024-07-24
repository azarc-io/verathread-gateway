package internal

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	pvtgraph "github.com/azarc-io/verathread-gateway/internal/gql/graph/private"
	pvtresolvers "github.com/azarc-io/verathread-gateway/internal/gql/graph/private/resolvers"
	pubgraph "github.com/azarc-io/verathread-gateway/internal/gql/graph/public"
	pubresolvers "github.com/azarc-io/verathread-gateway/internal/gql/graph/public/resolvers"
	middleware2 "github.com/azarc-io/verathread-gateway/internal/middleware"
	"github.com/azarc-io/verathread-gateway/internal/service"
	graphqluc "github.com/azarc-io/verathread-next-common/usecase/graphql"
	"github.com/erni27/imcache"

	apptypes "github.com/azarc-io/verathread-gateway/internal/types"

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
		log        zerolog.Logger // default logger
		opts       *apptypes.APIGatewayOptions
		httpClient *http.Client
		ready      bool
		is         apptypes.InternalService
		publicAPI  graphqluc.GraphQLUseCase
		privateAPI graphqluc.GraphQLUseCase
		proxy      *proxy
	}
)

/************************************************************************/
/* LIFECYCLE
/************************************************************************/

func (d *Domain) PreStart() error {
	// healthz
	healthz.Register("gateway", time.Second*1, func() error {
		if !d.ready {
			return apptypes.ErrGatewayNotReady
		}
		return nil
	})

	// create service to handle inbound requests
	d.is = service.NewService(d.opts, d.log)

	// register the application gateways own graphql endpoints
	if err := d.registerGqlAPI(); err != nil {
		return err
	}

	// register the shell app route
	d.registerShellAppRoute()

	// resource proxy
	d.registerProxyRouter()

	return nil
}

func (d *Domain) PostStart() error {
	// watches for leadership changed events in order to prevent contention on rebuilding of the shell's configuration
	// data in the cache
	d.opts.RedisUseCase.SubscribeToElectionEvents(func(onPromote <-chan time.Time, onDemote <-chan time.Time) {
		for {
			select {
			case <-onPromote:
				go func() {
					if err := d.is.Watch(); err != nil {
						panic(err)
					}
				}()
			case <-onDemote:
				if err := d.is.UnWatch(); err != nil {
					panic(err)
				}
			}
		}
	})

	// flag service is ready so health starts reporting ok status
	d.ready = true

	return nil
}

func (d *Domain) PreStop() error {
	// TODO mark unavailable so services return a service is going down error
	return nil
}

/************************************************************************/
/* SHELL APP
/************************************************************************/

// registerShellAppRoute serves up the shell app, supports 2 modes
// Proxy: If web proxy is set in your config then will proxy that url
// WebDir: If web director is set and proxy is not set then will serve static files from the web dir
func (d *Domain) registerShellAppRoute() {
	e := d.opts.PublicHTTPUseCase.Server()

	if d.opts.Config.WebProxy != "" {
		d.log.Info().Msgf("serving files from: %s", d.opts.Config.WebProxy)

		_url, err := url.Parse(d.opts.Config.WebProxy)
		if err != nil {
			panic(err)
		}

		tgt := &apptypes.ProxyTarget{
			ID:           "gateway",
			Name:         "shell",
			WebURL:       _url,
			APIURL:       _url,
			Meta:         nil,
			RegexRewrite: map[*regexp.Regexp]string{},
		}

		grp := e.Group("")
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

				if err := d.proxy.rewriteURL(tgt.RegexRewrite, req); err != nil {
					return err
				}

				c.Set(apptypes.TargetURLKey, tgt.APIURL)
				c.Set(apptypes.AppNameKey, tgt.Name)

				// Proxy
				switch {
				case c.IsWebSocket():
					d.proxy.proxyRaw(tgt, c).ServeHTTP(res, req)
				case req.Header.Get(echo.HeaderAccept) == "text/event-stream":
				default:
					d.proxy.proxyHTTP(tgt, c).ServeHTTP(res, req)
				}

				return nil
			}
		})
	} else if d.opts.Config.WebDir != "" {
		d.log.Info().Msgf("serving files from: %s", d.opts.Config.WebDir)

		g1 := e.Group("")
		g1.Use(
			func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					c.Response().Header().Set("Cache-Control", "max-age=31536000")
					return next(c)
				}
			},
			middleware.GzipWithConfig(middleware.GzipConfig{
				Skipper: func(c echo.Context) bool {
					return !strings.Contains(c.Path(), ".html")
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
/* APP PROXIES
/************************************************************************/

// registerProxyRouter registers the routes that are responsible for proxying both api requests and fetching
// modules for apps
func (d *Domain) registerProxyRouter() {
	// routes graph requests to an app by its service name, the app must have registered itself in advance
	grp1 := d.opts.PublicHTTPUseCase.Server().Group("/app/:appId/graphql")
	grp1.Use(middleware2.ACAOHeaderOverwriteMiddleware)
	grp1.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			app := c.Param("appId")

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

			if tgt, ok := d.is.GetProxyTarget(app); ok {
				if err := d.proxy.rewriteURL(tgt.RegexRewrite, req); err != nil {
					return err
				}

				d.log.Debug().Msgf("found gql proxy target <%s>:<%s>", app, tgt.APIURL)

				c.Set(apptypes.TargetURLKey, tgt.APIURL)
				c.Set(apptypes.AppNameKey, tgt.Name)

				// Proxy
				switch {
				case c.IsWebSocket():
					log.Debug().Msgf("proxy gql socket to %s%s", tgt.APIURL, req.URL)
					d.proxy.proxyRaw(tgt, c).ServeHTTP(res, req)
				case req.Header.Get(echo.HeaderAccept) == "text/event-stream":
					// TODO SSE
				default:
					log.Info().Msgf("proxy gql http   to %s%s", tgt.APIURL, req.URL)
					d.proxy.proxyHTTP(tgt, c).ServeHTTP(res, req)
				}
			} else {
				d.log.Warn().Msgf("no proxy target found for <%s>", app)
			}

			return nil
		}
	})

	// routes loading of web modules by service name, the app must have registered itself in advance
	grp2 := d.opts.PublicHTTPUseCase.Server().Group("/app/:appId")
	grp2.Use(middleware2.ACAOHeaderOverwriteMiddleware)
	grp2.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			app := c.Param("appId")

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

			if tgt, ok := d.is.GetProxyTarget(app); ok {
				if err := d.proxy.rewriteURL(tgt.RegexRewrite, req); err != nil {
					return err
				}

				d.log.Debug().Msgf("found web proxy target <%s>:<%s>", app, tgt.WebURL)

				c.Set(apptypes.TargetURLKey, tgt.WebURL)
				c.Set(apptypes.AppNameKey, tgt.Name)

				// Proxy
				switch {
				case c.IsWebSocket():
					log.Info().Msgf("proxy socket to %s%s", tgt.WebURL, req.URL)
					d.proxy.proxyRaw(tgt, c).ServeHTTP(res, req)
				case req.Header.Get(echo.HeaderAccept) == "text/event-stream":
					// TODO SSE
				default:
					log.Info().Msgf("proxy http   to %s%s", tgt.WebURL, req.URL)
					d.proxy.proxyHTTP(tgt, c).ServeHTTP(res, req)
				}
			} else {
				d.log.Warn().Msgf("no proxy target found for <%s>", app)
			}

			return nil
		}
	})
}

/************************************************************************/
/* API
/************************************************************************/

// registerGqlAPI registers graphql api handler
func (d *Domain) registerGqlAPI() error {
	// public api
	d.publicAPI = graphqluc.NewGraphQLUseCase(
		graphqluc.WithLogger(d.log),
		graphqluc.WithHTTPUseCase(d.opts.PublicHTTPUseCase),
		// graphqluc.WithServiceName(d.opts.ServiceName),
		graphqluc.WithExecutableSchema(pubgraph.NewExecutableSchema(pubgraph.Config{
			Resolvers: &pubresolvers.Resolver{
				Opts:            d.opts,
				InternalService: d.is,
			},
		})),
	)

	// private api
	d.privateAPI = graphqluc.NewGraphQLUseCase(
		graphqluc.WithLogger(d.log),
		graphqluc.WithHTTPUseCase(d.opts.PrivateHTTPUseCase),
		// graphqluc.WithServiceName(d.opts.ServiceName),
		graphqluc.WithExecutableSchema(pvtgraph.NewExecutableSchema(pvtgraph.Config{
			Resolvers: &pvtresolvers.Resolver{
				Opts:            d.opts,
				InternalService: d.is,
			},
		})),
	)

	return nil
}

/************************************************************************/
/* FACTORY
/************************************************************************/

func NewGateway(opts ...apptypes.APIGatewayOption) *Domain {
	l := log.With().Str("app", "gateway").Logger()
	g := &Domain{
		log:        l,
		opts:       &apptypes.APIGatewayOptions{},
		httpClient: http.DefaultClient,
		proxy: &proxy{
			log: l,
			httpProxyCache: imcache.NewSharded[string, *httputil.ReverseProxy](apptypes.CacheShards, imcache.DefaultStringHasher64{},
				imcache.WithCleanerOption[string, *httputil.ReverseProxy](apptypes.CacheCleanupFreq),
				imcache.WithEvictionCallbackOption[string, *httputil.ReverseProxy](func(key string, val *httputil.ReverseProxy, reason imcache.EvictionReason) {
					if reason == imcache.EvictionReasonExpired {
						log.Info().Str("uri", key).Msgf("http proxy evicted from cache")
					}
				}),
			),
		},
	}

	for _, opt := range opts {
		opt(g.opts)
	}

	// has to be set here so that options are applied first
	// this is the list of file names that should be scanned for tokens
	g.proxy.filesToScan = g.opts.Config.AssetsToScan

	return g
}
