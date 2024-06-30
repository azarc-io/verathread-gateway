package internal

import (
	"context"
	"encoding/json"
	"github.com/azarc-io/verathread-gateway/internal/api"
	"github.com/azarc-io/verathread-gateway/internal/cache"
	"github.com/azarc-io/verathread-gateway/internal/config"
	"github.com/azarc-io/verathread-next-common/common/app"
	"github.com/azarc-io/verathread-next-common/common/genericdb"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/rs/zerolog"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/datasource/httpclient"
)

type (
	RegistrationActor struct {
		cache  *cache.ProjectCache
		opts   *config.APIGatewayOptions
		log    zerolog.Logger
		db     genericdb.GenericDb
		client dapr.Client
		is     api.InternalService
	}
)

func (a *RegistrationActor) Register(ctx context.Context, v *common.InvocationEvent) (*common.Content, error) {
	var req app.RegisterAppInput
	if err := json.Unmarshal(v.Data, &req); err != nil {
		return nil, err
	}

	rsp, err := a.is.RegisterApp(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &common.Content{Data: rsp.MustMarshal(), ContentType: httpclient.ContentTypeJSON}, nil
}

func (a *RegistrationActor) KeepAlive(ctx context.Context, v *common.InvocationEvent) (*common.Content, error) {
	var req app.KeepAliveAppInput
	if err := json.Unmarshal(v.Data, &req); err != nil {
		return nil, err
	}

	rsp, err := a.is.KeepAlive(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &common.Content{
		Data:        rsp.MustMarshal(),
		ContentType: httpclient.ContentTypeJSON,
	}, nil
}

func (a *RegistrationActor) RegisterHandlers() error {
	a.log.Info().Msgf("creating app registration handler")
	if err := a.opts.DaprUseCase.Service().AddServiceInvocationHandler("/app/register", a.Register); err != nil {
		return err
	}

	a.log.Info().Msgf("creating app keep alive handler")
	if err := a.opts.DaprUseCase.Service().AddServiceInvocationHandler("/app/keep-alive", a.KeepAlive); err != nil {
		return err
	}

	return nil
}

func newAppRegistrationHandler(cache *cache.ProjectCache, opts *config.APIGatewayOptions, log zerolog.Logger, is api.InternalService) *RegistrationActor {
	return &RegistrationActor{
		cache:  cache,
		opts:   opts,
		client: opts.DaprUseCase.Client(),
		log:    log,
		db:     opts.MongoUseCase.GenericClient(),
		is:     is,
	}
}
