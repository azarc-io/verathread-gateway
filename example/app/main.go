package main

import (
	"github.com/azarc-io/verathread-next-common/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"

	appuc "github.com/azarc-io/verathread-next-common/usecase/app"

	"github.com/azarc-io/verathread-next-common/common/app"
	"github.com/azarc-io/verathread-next-common/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	Domain = "example"
)

type (
	example struct {
		log zerolog.Logger // active logger
		svc *service.AppService
		cfg *Config
	}

	Config struct {
		service.Config `yaml:",inline"`
		Registration   *Registration `yaml:"registration"`
		WebDir         string        `yaml:"webDir"`
	}

	Registration struct {
		WebBaseURL string `yaml:"web_url"`
		APIBaseURL string `yaml:"api_url"`
		GatewayURL string `yaml:"gateway_url"`
	}
)

// main entry point, the gateway should be running locally before you start this example app
// it will register itself with the gateway and should cause the navigation in the shell to update
// you should then be able to navigate to the examples micro front end
func main() {
	// configuration
	var cfg Config
	util.PanicIfErr(service.LoadConfig(&cfg))

	a := &example{
		log: log.With().Str("app", "example").Logger(),
		cfg: &cfg,
	}

	a.svc = service.NewAppService(
		service.WithConfig(&cfg.Config),
		service.WithBeforeStart(func(svc *service.AppService) error {
			return a.registerWebHandler()
		}),
		service.WithAfterStart(func(svc *service.AppService) error {
			return a.registerApp()
		}),
		service.WithBeforeStop(func() error {
			return nil
		}),
	)
	util.PanicIfErr(a.svc.Run())
}

// registerWebHandler serves up the web app
func (ex *example) registerWebHandler() error {
	e := ex.svc.PublicHTTP().Server()

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
			Root:  ex.cfg.WebDir,
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

	return nil
}

func (ex *example) registerApp() error {
	// initialize the app registration use case
	reg := appuc.NewAppUseCase(
		appuc.WithGatewayUrl(ex.cfg.Registration.GatewayURL),
		appuc.WithLogger(ex.log),
		appuc.WithAppInfo(app.RegisterAppInput{
			Id:              "gateway-example-1",
			Name:            "gateway-example", // the gateway will use this name to proxy e.g. /module/user/*
			Package:         "vth:azarc:gateway-example",
			Version:         "1.0.0", // TODO inject from ci and use here
			ApiUrl:          ex.cfg.Registration.APIBaseURL,
			RemoteEntryFile: "remoteEntry.js", // if proxy is true then don't need url here
			WebUrl:          ex.cfg.Registration.WebBaseURL,
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
						ExposedModule: "./Counter",
						ModuleName:    "example",
						Path:          "/example",
						Outlet:        "",
					},
					AuthRequired: true,
					Hidden:       false,
					Proxy:        true,
					Category:     app.RegisterAppCategorySetting,
					Children:     []app.RegisterChildAppNavigationInput{},
					Icon:         "",
				},
			},
		}),
	)
	return reg.Start()
}
