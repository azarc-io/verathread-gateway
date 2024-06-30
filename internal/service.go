package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/azarc-io/verathread-gateway/internal/api"
	"github.com/azarc-io/verathread-gateway/internal/cache"
	"github.com/azarc-io/verathread-gateway/internal/config"
	"github.com/azarc-io/verathread-gateway/internal/gql/graph/model"
	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
	apputil "github.com/azarc-io/verathread-gateway/internal/util"
	"github.com/azarc-io/verathread-next-common/common/app"
	"github.com/azarc-io/verathread-next-common/common/genericdb"
	util2 "github.com/azarc-io/verathread-next-common/util"
	hashutil "github.com/azarc-io/verathread-next-common/util/hash"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	service struct {
		log   zerolog.Logger
		opts  *config.APIGatewayOptions
		db    *mongo.Collection
		cache *cache.ProjectCache
	}
)

func (s service) RegisterApp(ctx context.Context, req *app.RegisterAppInput) (*app.RegisterAppOutput, error) {
	id := hashutil.UuidFromString(req.Package)
	count, err := s.db.CountDocuments(ctx, bson.M{
		"_id": id,
	})

	if err != nil && !errors.Is(err, genericdb.ErrRecordNotFound) {
		s.log.Error().Err(err).Msgf("failed to check for existence of ent")
		return nil, err
	}

	ent := &apptypes.App{
		Name:        req.Name,
		Package:     req.Package,
		Version:     req.Version,
		ApiUrl:      req.ApiUrl,
		ApiWsUrl:    req.ApiWsUrl,
		BaseUrl:     req.BaseUrl,
		RemoteEntry: req.RemoteEntryFile,
		ProxyApi:    req.ProxyApi,
		Proxy:       req.Proxy,
		Navigation:  []*apptypes.Navigation{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Now(),
		Adopted:     util2.Ptr(true),
		Available:   util2.Ptr(true),
	}

	if count == 0 {
		ent.CreatedAt = time.Now()
	}

	for _, navigation := range req.Navigation {
		n := &apptypes.Navigation{
			Id: hashutil.UuidFromString(req.Package),
		}

		apputil.MapNavigationToNavigationInput(n, navigation)

		ent.Navigation = append(ent.Navigation, n)

		if navigation.Proxy {
			n.RemoteEntryRewriteRegEx = map[string]string{
				"/module/*/*": "/$2",
			}
			n.RemoteEntry = util2.Ptr(
				fmt.Sprintf("%s/module/%s/remoteEntry.js", "", n.Id)) //a.opts.Config.GatewayBaseUrl
		} else {
			n.RemoteEntry = util2.Ptr(fmt.Sprintf("%s/%s", req.BaseUrl, req.RemoteEntryFile))
		}
	}

	if req.Slot1 != nil {
		ent.Slot1 = apputil.MapRegisterSlotToEntity(req.Slot1)
	}

	if req.Slot2 != nil {
		ent.Slot2 = apputil.MapRegisterSlotToEntity(req.Slot2)
	}

	if req.Slot3 != nil {
		ent.Slot3 = apputil.MapRegisterSlotToEntity(req.Slot3)
	}

	if _, err := s.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": ent,
	}, options.Update().SetUpsert(true)); err != nil {
		s.log.Error().Err(err).Msgf("failed to register app")
		return nil, err
	}

	ent.Id = id

	s.log.Info().Str("pkg", req.Package).Msgf("registered app")
	s.cache.Add(ent, time.Now().Add(time.Second*4))

	// TODO decide if we need to publish these events or not
	//if count == 0 {
	//	if err := a.client.PublishEvent(ctx, "vth-ent-stream", "ent.v1.registered", ent); err != nil {
	//		a.log.Error().Err(err).Msgf("failed to publish event")
	//	}
	//} else {
	//	if err := a.client.PublishEvent(ctx, "vth-ent-stream", "ent.v1.updated", ent); err != nil {
	//		a.log.Error().Err(err).Msgf("failed to publish event")
	//	}
	//}

	return &app.RegisterAppOutput{Id: id}, nil
}

func (s service) KeepAlive(ctx context.Context, req *app.KeepAliveAppInput) (*app.KeepAliveAppOutput, error) {
	if _, ok := s.cache.Get(req.Pkg); ok {
		s.cache.ResetExpiryOf(req.Pkg, time.Second*4)
		rsp := &app.KeepAliveAppOutput{
			RegistrationRequired: false,
			Ok:                   true,
		}
		return rsp, nil
	}

	return &app.KeepAliveAppOutput{
		RegistrationRequired: true,
		Ok:                   false,
	}, nil
}

func (s service) GetAppConfiguration(ctx context.Context, tenant string) (*model.ShellConfiguration, error) {
	cur, err := s.opts.MongoUseCase.Collection("app").Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	var apps []*apptypes.App

	if err := cur.All(context.Background(), &apps); err != nil {
		return nil, err
	}

	return apputil.MapAppsToNavigation(apps), nil
}

func NewService(opts *config.APIGatewayOptions, log zerolog.Logger, cache *cache.ProjectCache) api.InternalService {
	return &service{
		log:   log,
		opts:  opts,
		db:    opts.MongoUseCase.Collection("app"),
		cache: cache,
	}
}
