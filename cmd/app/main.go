package main

import (
	"context"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/azarc-io/verathread-gateway/internal"
	mongouc "github.com/azarc-io/verathread-next-common/usecase/mongo"

	"github.com/azarc-io/verathread-gateway/internal/config"
	authzuc "github.com/azarc-io/verathread-next-common/usecase/authz"
	dapruc "github.com/azarc-io/verathread-next-common/usecase/dapr"
	devuc "github.com/azarc-io/verathread-next-common/usecase/dev"
	luc "github.com/azarc-io/verathread-next-common/usecase/logging"
	wardenuc "github.com/azarc-io/verathread-next-common/usecase/warden"
	"github.com/azarc-io/verathread-next-common/util"
	cfgutil "github.com/azarc-io/verathread-next-common/util/config"
	signals "github.com/azarc-io/verathread-next-common/util/signal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	app struct {
		ctx     context.Context      // global context
		cancel  context.CancelFunc   // called on service shutdown, cancels global context
		log     zerolog.Logger       // active logger
		cfg     *config.Config       // composed config
		gateway *internal.Domain     // main api entry point
		auth    authzuc.AuthZUseCase // auth helper
		// tracing tracinguc.TracingUseCase
		warden wardenuc.ClusterWardenUseCase
		dapr   dapruc.DaprUseCase
		mongo  mongouc.MongoUseCase
	}
)

func main() {
	// start times
	startAt := time.Now()

	// global contexts
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		if err := recover(); err != nil {
			log.Error().Err(err.(error)).Msgf("service did not shutdown cleanly\n%s", debug.Stack())
			// cancel global context
			cancel()
			os.Exit(1)
		}
	}()

	// configuration
	var cfg *config.Config
	if err := cfgutil.LoadContextual(
		os.Getenv("CONFIG_DIR"),
		os.Getenv("BASE_CONTEXT"),
		cfgutil.SubContextsFromEnv("CONFIG_CONTEXTS"),
		&cfg,
		".env",
	); err != nil {
		panic(err)
	}

	// configure logging
	luc.NewLoggingUseCase(
		luc.WithLevel(cfg.Logger.Level),
		luc.WithMode(cfg.Logger.Mode),
	)

	a := &app{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
		log:    log.With().Str("app", "standalone").Logger(),
	}

	// Tracing use case
	// TODO re-enable once metrics and traces are enabled for the metrics agent
	// a.tracing = tracinguc.NewZipkinTracer(
	//	tracinguc.WithServiceName(cfg.Name),
	//	tracinguc.WithUrl("http://localhost:9411/api/v2/spans"),
	// )
	// util.PanicIfErr(a.tracing.Start())

	// dapr use case, provided client and service with some helpers
	a.dapr = dapruc.NewDaprUseCase(
		dapruc.WithLogger(a.log),
		dapruc.WithConfig(a.cfg.Dapr),
		dapruc.WithServicePort(a.cfg.Gateway.HTTP.Port),
	)

	// auth use case
	a.auth = authzuc.NewAuthZUseCase(
		authzuc.WithConfig(a.cfg.Auth),
	)
	util.PanicIfErr(a.auth.Start())

	// mongo use case, start immediately in case some actors want to access the db early on
	a.mongo = mongouc.NewMongoUseCase(mongouc.WithConfig(cfg.Database))
	util.PanicIfErr(a.mongo.Start())

	// warden authorization use case
	a.warden = wardenuc.NewClusteredWardenUseCase(
		wardenuc.WithCommonAuthUseCase(a.auth),
	)

	a.gateway = internal.NewGateway(
		config.WithConfig(a.cfg.Gateway),
		config.WithServiceID(a.cfg.ID),
		config.WithAuthUseCase(a.auth),
		config.WithWardenUseCase(a.warden),
		config.WithDaprUseCase(a.dapr),
		config.WithMongoUseCase(a.mongo),
	)

	// initialize dev mode if enabled
	stoppedCh := make(chan struct{})
	<-a.initDevMode(ctx, stoppedCh)

	// start the client after dapr side-car has started
	util.PanicIfErr(a.dapr.StartClient())

	// because we are sharing the actor system we need to init the gateway, so it can register its actors
	util.PanicIfErr(a.gateway.Init())

	// start dapr service after domain is started
	util.PanicIfErr(a.dapr.StartService())

	// finally start the gateway
	util.PanicIfErr(a.gateway.Start())

	// log out start up time
	a.log.Info().Msgf("server started in %s", time.Since(startAt))

	// wait for shutdown signals
	<-signals.SetupSignalHandler()

	// gracefully shutdown services
	a.log.Info().Msgf("shutting down gracefully after %s", time.Since(startAt))

	// stop the gateway gracefully
	if err := a.gateway.Stop(); err != nil {
		a.log.Error().Err(err).Msgf("encountered error while gracefully shutting down the gateway")
	}

	// cancel global context
	cancel()
	<-stoppedCh
}

func (a *app) initDevMode(ctx context.Context, ch chan struct{}) chan struct{} {
	uc := devuc.NewDevelopmentUseCase(devuc.WithConfig(a.cfg.Development))
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

var GcTickDuration = time.Second * 10

func init() {
	// try to free up any memory periodically
	go func() {
		tick := time.Tick(GcTickDuration)
		for range tick {
			runtime.GC()
			debug.FreeOSMemory()
		}
	}()
}
