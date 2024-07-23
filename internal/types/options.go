package types

import (
	"context"

	authzuc "github.com/azarc-io/verathread-next-common/usecase/authz"
	httpuc "github.com/azarc-io/verathread-next-common/usecase/http"
	mongouc "github.com/azarc-io/verathread-next-common/usecase/mongo"
	natsuc "github.com/azarc-io/verathread-next-common/usecase/nats"
	redisuc "github.com/azarc-io/verathread-next-common/usecase/redis"
	wardenuc "github.com/azarc-io/verathread-next-common/usecase/warden"
)

type (
	APIGatewayOptions struct {
		ServiceID          string
		Config             *APIGatewayConfig
		AuthUseCase        authzuc.AuthZUseCase
		WardenUseCase      wardenuc.ClusterWardenUseCase
		MongoUseCase       mongouc.MongoUseCase
		PublicHTTPUseCase  httpuc.HttpUseCase
		PrivateHTTPUseCase httpuc.HttpUseCase
		ServiceName        string
		RedisUseCase       redisuc.RedisUseCase
		Context            context.Context
		NatsUseCase        natsuc.NatsUseCase
	}

	APIGatewayOption func(o *APIGatewayOptions)

	APIGatewayConfig struct {
		WebDir        string                 `yaml:"web_dir"`
		WebProxy      string                 `yaml:"web_proxy"`
		HTTP          *httpuc.ConfigBindHttp `yaml:"http"`
		Services      map[string]Service     `yaml:"services"`
		BackofficeOrg string                 `yaml:"backoffice_org"`
		AssetsToScan  []string               `yaml:"assets_to_scan"`
	}

	Service struct {
		Gql   string `yaml:"gql"`
		GqlWs string `yaml:"gql_ws"`
	}
)

func WithContext(ctx context.Context) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.Context = ctx
	}
}

func WithServiceID(id string) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.ServiceID = id
	}
}

func WithConfig(cfg *APIGatewayConfig) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.Config = cfg
	}
}

func WithAuthUseCase(auc authzuc.AuthZUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.AuthUseCase = auc
	}
}

func WithNatsUseCase(nuc natsuc.NatsUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.NatsUseCase = nuc
	}
}

func WithWardenUseCase(wuc wardenuc.ClusterWardenUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.WardenUseCase = wuc
	}
}

func WithMongoUseCase(muc mongouc.MongoUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.MongoUseCase = muc
	}
}

func WithPublicHTTPUseCase(huc httpuc.HttpUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.PublicHTTPUseCase = huc
	}
}

func WithPrivateHTTPUseCase(huc httpuc.HttpUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.PrivateHTTPUseCase = huc
	}
}

func WithRedisUseCase(ruc redisuc.RedisUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.RedisUseCase = ruc
	}
}

func WithServiceName(name string) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.ServiceName = name
	}
}
