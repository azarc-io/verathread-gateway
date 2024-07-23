package types

import (
	"context"

	"github.com/azarc-io/verathread-gateway/internal/gql/graph/common/model"
)

type (
	AppDomain interface {
		Start() error
		Stop() error
		PostStart() error
	}

	InternalService interface {
		GetAppConfiguration(ctx context.Context, tenant string) (*model.ShellConfiguration, error)
		RegisterApp(ctx context.Context, req *model.RegisterAppInput) (*model.RegisterAppOutput, error)
		KeepAlive(ctx context.Context, req *model.KeepAliveAppInput) (*model.KeepAliveAppOutput, error)
		GetProxyTarget(app string) (*ProxyTarget, bool)
		Watch() error
		UnWatch() error
	}
)
