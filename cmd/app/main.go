package main

import (
	"os"
	"runtime"
	"runtime/debug"
	"time"

	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
	"github.com/azarc-io/verathread-next-common/service"
	wardenuc "github.com/azarc-io/verathread-next-common/usecase/warden"
	"github.com/azarc-io/verathread-next-common/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/azarc-io/verathread-gateway/internal"
)

type (
	app struct {
		log     zerolog.Logger   // active logger
		gateway *internal.Domain // main api entry point
		warden  wardenuc.ClusterWardenUseCase
		svc     *service.AppService
	}

	Config struct {
		service.Config `yaml:",inline"`
		Name           string                     `yaml:"name"`
		ID             string                     `yaml:"id"`
		DataDir        string                     `yaml:"data_dir"`
		Gateway        *apptypes.APIGatewayConfig `yaml:"gateway"`
	}
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Err(err.(error)).Msgf("service did not startup cleanly\n%s", debug.Stack())
			// cancel global context
			os.Exit(1)
		}
	}()

	// configuration
	var cfg Config
	util.PanicIfErr(service.LoadConfig(&cfg))

	a := &app{
		log: log.With().Str("app", "account").Logger(),
	}

	a.svc = service.NewAppService(
		service.WithConfig(&cfg.Config),
		service.WithBeforeStart(func(svc *service.AppService) error {
			a.gateway = internal.NewGateway(
				apptypes.WithConfig(cfg.Gateway),
				apptypes.WithServiceID(cfg.ID),
				apptypes.WithServiceName(cfg.Name),
				apptypes.WithAuthUseCase(svc.Auth()),
				apptypes.WithWardenUseCase(a.warden),
				apptypes.WithMongoUseCase(svc.Mongo()),
				apptypes.WithPrivateHTTPUseCase(svc.PrivateHTTP()),
				apptypes.WithPublicHTTPUseCase(svc.PublicHTTP()),
				apptypes.WithRedisUseCase(svc.Redis()),
				apptypes.WithContext(svc.Context()),
				apptypes.WithNatsUseCase(svc.Nats()),
			)
			return a.gateway.PreStart()
		}),
		service.WithAfterStart(func(svc *service.AppService) error {
			return a.gateway.PostStart()
		}),
		service.WithBeforeStop(func() error {
			return a.gateway.PreStop()
		}),
	)
	util.PanicIfErr(a.svc.Run())
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
