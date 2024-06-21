package config

import (
	authzuc "github.com/azarc-io/verathread-next-common/usecase/authz"
	dapruc "github.com/azarc-io/verathread-next-common/usecase/dapr"
	httpuc "github.com/azarc-io/verathread-next-common/usecase/http"
	wardenuc "github.com/azarc-io/verathread-next-common/usecase/warden"
)

type (
	APIGatewayOptions struct {
		ServiceID     string
		Config        *APIGatewayConfig
		AuthUseCase   authzuc.AuthZUseCase
		WardenUseCase wardenuc.ClusterWardenUseCase
		DaprUseCase   dapruc.DaprUseCase
	}

	APIGatewayOption func(o *APIGatewayOptions)

	APIGatewayConfig struct {
		WebDir        string                 `yaml:"web_dir"`
		WebProxy      string                 `yaml:"web_proxy"`
		HTTP          *httpuc.ConfigBindHttp `yaml:"http"`
		Services      map[string]Service     `yaml:"services"`
		BackofficeOrg string                 `yaml:"backoffice_org"`
	}

	Service struct {
		Gql   string `yaml:"gql"`
		GqlWs string `yaml:"gql_ws"`
	}
)

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

func WithWardenUseCase(wuc wardenuc.ClusterWardenUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.WardenUseCase = wuc
	}
}

func WithDaprUseCase(duc dapruc.DaprUseCase) APIGatewayOption {
	return func(o *APIGatewayOptions) {
		o.DaprUseCase = duc
	}
}
