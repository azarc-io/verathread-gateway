package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	appuc "github.com/azarc-io/verathread-next-common/usecase/app"

	"github.com/azarc-io/verathread-next-common/common/app"
	dapruc "github.com/azarc-io/verathread-next-common/usecase/dapr"
	devuc "github.com/azarc-io/verathread-next-common/usecase/dev"
	httpuc "github.com/azarc-io/verathread-next-common/usecase/http"
	luc "github.com/azarc-io/verathread-next-common/usecase/logging"
	"github.com/azarc-io/verathread-next-common/util"
	cfgutil "github.com/azarc-io/verathread-next-common/util/config"
	signals "github.com/azarc-io/verathread-next-common/util/signal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

const (
	Domain = "example"
)

type (
	Config struct {
		Dapr         *dapruc.Config         `yaml:"dapr"`
		Development  *devuc.Config          `yaml:"development"`
		Registration *Registration          `yaml:"registration"`
		HTTP         *httpuc.ConfigBindHttp `yaml:"http"`
		WebDir       string                 `yaml:"webDir"`
	}

	Registration struct {
		WebBaseURL string `yaml:"webBaseUrl"`
		BaseWsURL  string `yaml:"baseWsUrl"`
	}
)

// main entry point, the gateway should be running locally before you start this example app
// it will register itself with the gateway and should cause the navigation in the shell to update
// you should then be able to navigate to the examples micro front end
func main() {
	luc.NewLoggingUseCase(luc.WithMode(luc.LoggingDevMode), luc.WithLevel("info"))
	l := log.With().Str("app", "example").Logger()

	// global contexts
	ctx, cancel := context.WithCancel(context.Background())

	// configuration
	var cfg *Config
	if err := cfgutil.LoadContextual(
		os.Getenv("CONFIG_DIR"),
		os.Getenv("BASE_CONTEXT"),
		cfgutil.SubContextsFromEnv("CONFIG_CONTEXTS"),
		&cfg,
		".env",
	); err != nil {
		panic(err)
	}

	// dapr use case, provided client and service with some helpers
	dapr := dapruc.NewDaprUseCase(
		dapruc.WithLogger(l),
		dapruc.WithConfig(cfg.Dapr),
		dapruc.WithServicePort(cfg.HTTP.Port),
	)

	// create http server
	huc := httpuc.NewHttpUseCase(
		httpuc.WithHttpConfig(cfg.HTTP),
		httpuc.WithLogger(l),
	)

	registerWebHandler(huc, cfg)

	// route all unhandled requests to echo
	mux := dapr.Mux()
	mux.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		huc.Server().ServeHTTP(writer, request)
	})

	// initialize dev mode if enabled
	stoppedCh := make(chan struct{})
	<-initDevMode(ctx, cfg.Development, stoppedCh)

	// initialize the app registration use case
	app := appuc.NewAppUseCase(
		appuc.WithGatewayUrl("http://dev.cluster.local/grqphql"),
		appuc.WithLogger(l),
		appuc.WithAppInfo(app.RegisterAppInput{
			Name:            "gateway-example", // the gateway will use this name to proxy e.g. /module/user/*
			Package:         "vth:azarc:gateway:example",
			Version:         "1.0.0", // TODO inject from ci and use here
			ApiUrl:          fmt.Sprintf("http://%s/graphql", net.JoinHostPort(cfg.HTTP.Address, strconv.Itoa(cfg.HTTP.Port))),
			RemoteEntryFile: "remoteEntry.js", // if proxy is true then don't need url here
			WebUrl:          fmt.Sprintf("%s/app/%s", cfg.Registration.WebBaseURL, Domain),
			Proxy:           false,
			Slot1: app.RegisterAppSlot{
				Description:  "Slot 1 module has no path so it must be a drop down",
				AuthRequired: false,
				Module: app.RegisterAppSlotModule{
					ExposedModule: "./AppSlot1Module",
					ModuleName:    "AppSlot1Module",
				},
			},
			Slot2: app.RegisterAppSlot{
				Description:  "Slot 2 module has no path so it must be a drop down",
				AuthRequired: false,
				Module: app.RegisterAppSlotModule{
					ExposedModule: "./AppSlot2Module",
					ModuleName:    "AppSlot2Module",
				},
			},
			Slot3: app.RegisterAppSlot{
				Description:  "Slot 3 module has a path so it just a shortcut to a navigable path",
				AuthRequired: false,
				Module: app.RegisterAppSlotModule{
					ExposedModule: "./AppSlot3Module",
					ModuleName:    "AppSlot3Module",
					Path:          "/rune",
				},
			},
			Navigation: []app.RegisterAppNavigationInput{
				{
					Title:    "Example App Root",
					SubTitle: "Example root entry",
					Module: app.RegisterAppModule{
						ExposedModule: "./AppModule",
						ModuleName:    "ExampleModule",
						Path:          "/rune",
					},
					AuthRequired: true,
					Hidden:       false,
					Proxy:        true,
					Category:     app.RegisterAppCategorySetting,
					Children: []app.RegisterChildAppNavigationInput{
						{
							Title:        "Example App Child 1",
							SubTitle:     "Example child entry",
							AuthRequired: true,
							Path:         "example1",
						},
					},
				},
			},
		}),
	)

	l.Info().Msgf("initializing client")
	util.PanicIfErr(dapr.StartClient())

	l.Info().Msgf("starting registration loop")
	util.PanicIfErr(app.Start())

	l.Info().Msgf("initializing service")
	util.PanicIfErr(dapr.StartService())

	l.Info().Msgf("service started")

	// wait for shutdown signals
	<-signals.SetupSignalHandler()
	cancel()
}

func initDevMode(ctx context.Context, cfg *devuc.Config, ch chan struct{}) chan struct{} {
	uc := devuc.NewDevelopmentUseCase(devuc.WithConfig(cfg))
	ready := make(chan struct{})
	go func() {
		_, stop := uc.Start()
		close(ready)
		<-ctx.Done()
		if err := stop(); err != nil {
			log.Error().Err(err).Msgf("error while shuttind down dapr")
		}
		log.Info().Msgf("dapr stopped gracefully")
		close(ch)
	}()

	return ready
}

// registerWebHandler serves up the web app
func registerWebHandler(huc httpuc.HttpUseCase, cfg *Config) {
	e := huc.Server()

	g1 := e.Group("/module/" + Domain)
	g1.Use(
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				ct := c.Response().Header().Get(echo.HeaderContentType)
				return ct != "text/css" &&
					!strings.HasPrefix(ct, "application/javascript") &&
					!strings.HasPrefix(ct, "text/javascript")
			},
		}),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  cfg.WebDir,
			Index: "remoteEntry.json",
			HTML5: true,
			Skipper: func(e echo.Context) bool {
				return strings.HasPrefix(e.Path(), "/tmp") ||
					strings.HasPrefix(e.Path(), "/api") ||
					strings.HasPrefix(e.Path(), "/graphql") ||
					strings.HasPrefix(e.Path(), "/query")
			},
		}),
	)
}
