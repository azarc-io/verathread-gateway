package cache

import (
	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
	"github.com/erni27/imcache"
	"github.com/rs/zerolog"
	"time"
)

type (
	ProjectCache struct {
		cache *imcache.Sharded[string, *apptypes.App]
		log   zerolog.Logger
	}
)

func (c *ProjectCache) Add(app *apptypes.App, expiresAt time.Time) {
	c.cache.Set(app.Package, app, imcache.WithExpirationDate(expiresAt))
}

func (c *ProjectCache) ResetExpiryOf(pkg string, duration time.Duration) {
	newExpiry := time.Now().Add(duration)
	c.cache.ReplaceWithFunc(pkg, func(project *apptypes.App) *apptypes.App {
		return project
	}, imcache.WithExpirationDate(newExpiry))
}

func (c *ProjectCache) Get(pkg string) (*apptypes.App, bool) {
	return c.cache.Get(pkg)
}

func NewProjectCache(logger zerolog.Logger) *ProjectCache {
	return &ProjectCache{
		log: logger,
		cache: imcache.NewSharded[string, *apptypes.App](4, imcache.DefaultStringHasher64{},
			imcache.WithCleanerOption[string, *apptypes.App](time.Second*5),
			imcache.WithEvictionCallbackOption[string, *apptypes.App](func(key string, val *apptypes.App, reason imcache.EvictionReason) {
				if reason == imcache.EvictionReasonExpired {
					logger.Warn().
						Str("name", val.Name).
						Str("pkg", val.Package).
						Str("ver", val.Version).
						Str("reason", reason.String()).
						Msgf("app evicted from cache, service is down or not sending keep alive messages")

					//auc.System().Root.Send(
					//	auc.System().NewLocalPID(appapi.AppEventManagerActorName),
					//	&apptypes.AppUnavailable{App: val},
					//)
				}
			}),
		),
	}
}
