package api

import (
	"context"
	"github.com/azarc-io/verathread-gateway/internal/gql/graph/model"
	"github.com/azarc-io/verathread-next-common/common/app"
)

type (
	AppDomain interface {
		Start() error
		Stop() error
		PostStart() error
	}

	InternalService interface {
		GetAppConfiguration(ctx context.Context, tenant string) (*model.ShellConfiguration, error)
		RegisterApp(ctx context.Context, req *app.RegisterAppInput) (*app.RegisterAppOutput, error)
		KeepAlive(ctx context.Context, req *app.KeepAliveAppInput) (*app.KeepAliveAppOutput, error)
	}
)
