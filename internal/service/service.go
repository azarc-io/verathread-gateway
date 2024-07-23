package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/azarc-io/verathread-gateway/internal/gql/graph/common/model"
	"github.com/erni27/imcache"
	"github.com/redis/go-redis/v9"

	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
	apputil "github.com/azarc-io/verathread-gateway/internal/util"
	hashutil "github.com/azarc-io/verathread-next-common/util/hash"
	"github.com/rs/zerolog"
)

type (
	service struct {
		log  zerolog.Logger
		opts *apptypes.APIGatewayOptions
		sync.Mutex
		watchSub    *redis.PubSub
		targetCache *imcache.Sharded[string, *apptypes.ProxyTarget]
	}
)

/************************************************************************/
/* KEYSPACE NOTIFICATION
/************************************************************************/

// Watch runs only on a single instance in the cluster and watches the redis keyspace for notifications.
// - whenever a key has expired it will update the availability flag of the app and regenerate the cached configuration
func (s *service) Watch() error {
	rc := s.opts.RedisUseCase.Client()

	s.log.Info().Msgf("watching for application cache keyspace changes")

	if err := s.rebuildNavigation(); err != nil {
		s.log.Warn().Err(err).Msgf("could not rebuild navigation")
	}

	s.watchSub = rc.PSubscribe(s.opts.Context, apptypes.KeySpaceExpiryChannel)

	go func() {
		for {
			msg, err := s.watchSub.ReceiveMessage(s.opts.Context)
			if err != nil {
				s.log.Warn().Err(err).Msgf("failed to receive app exiry event from redis")
				continue
			}

			if strings.HasPrefix(msg.Payload, apptypes.KeepAliveKeySpacePrefix) {
				id := strings.Split(msg.Payload, ":")
				s.log.Info().Msgf("handling app event: %s => %s", msg.Channel, msg.Payload)

				if err := s.markApplicationAsUnavailable(id[2]); err != nil {
					s.log.Warn().Err(err).Msgf("could not mark app as unavailable on keyspace event")
					continue
				}

				if err := s.rebuildNavigation(); err != nil {
					s.log.Warn().Err(err).Msgf("could not rebuild navigation on keyspace event")
					continue
				}
			}
		}
	}()

	return nil
}

// UnWatch called when an instance looses leadership and stops watching for keyspace events
func (s *service) UnWatch() error {
	if s.watchSub != nil {
		return s.watchSub.Close()
	}

	return nil
}

/************************************************************************/
/* PROXY HELPERS
/************************************************************************/

// GetProxyTarget checks the local cache for a proxy targets configuration, if not found it loads the
// app from the cache and generates a proxy config that is then cached for a period of time
func (s *service) GetProxyTarget(appID string) (*apptypes.ProxyTarget, bool) {
	var (
		target *apptypes.ProxyTarget
		exists bool
		rc     = s.opts.RedisUseCase.Client()
		app    apptypes.App
	)

	if target, exists = s.targetCache.Get(appID); !exists {
		cmd := rc.HGet(s.opts.Context, "apps", appID)
		if cmd.Err() != nil {
			return nil, false
		}
		if err := cmd.Scan(&app); err != nil {
			s.log.Warn().Err(err).Msgf("fetched cached app but could not unmarshal the data")
			return nil, false
		}

		webURL, err := url.Parse(app.WebURL)
		if err != nil {
			s.log.Error().Err(err).Msgf("failed to parse base url for application")
			return nil, false
		}

		apiURL, err := url.Parse(app.APIURL)
		if err != nil {
			s.log.Error().Err(err).Msgf("failed to parse base url for application")
			return nil, false
		}

		target = &apptypes.ProxyTarget{
			ID:           app.ID,
			Name:         app.Name,
			WebURL:       webURL,
			APIURL:       apiURL,
			Meta:         map[string]interface{}{}, // TODO fill in auth etc.
			RegexRewrite: make(map[*regexp.Regexp]string),
		}

		if app.RemoteEntryRewriteRegEx != nil {
			for k, v := range apputil.RewriteRulesRegex(app.RemoteEntryRewriteRegEx) {
				target.RegexRewrite[k] = v
			}
		}

		s.targetCache.Set(appID, target, imcache.WithExpiration(apptypes.TargetCacheDuration))
		exists = true
	}

	return target, exists
}

/************************************************************************/
/* APP REGISTRATION
/************************************************************************/

// RegisterApp registers an application, invoked through the gql api and a record to the hash set in the cache
// and the record will not be removed if the app stops sending keep alive notifications, instead the app will
// be marked as unhealthy. An app is only removed when it's uninstalled through a user action or the cache is flushed
// and the app is also offline.
func (s *service) RegisterApp(ctx context.Context, req *model.RegisterAppInput) (*model.RegisterAppOutput, error) {
	var (
		ent       apptypes.App
		rc        = s.opts.RedisUseCase.Client()
		err       error
		appExists bool
		appKey    = req.Name
	)

	er := rc.HExists(ctx, "apps", appKey)
	if err = er.Err(); err != nil {
		s.log.Error().Str("package", req.Package).Err(err).Msgf("failed to retrieve check for cached app entry")
		return nil, fmt.Errorf("failed to retrieve check for cached app entry: %w", err)
	}
	appExists = er.Val()

	if appExists {
		gar := rc.HGet(ctx, "apps", appKey)
		if err = gar.Err(); err != nil {
			s.log.Error().Str("package", req.Package).Err(err).Msgf("failed to retrieve cached app entry")
			return nil, fmt.Errorf("failed to retrieve cached app entry: %w", err)
		}

		if err := gar.Scan(&ent); err != nil {
			s.log.Error().Str("package", req.Package).Err(err).Msgf("failed to scan cached app")
			return nil, fmt.Errorf("failed to scan cached app: %w", err)
		}
	} else {
		ent = apptypes.App{}
	}

	if !appExists {
		ent.CreatedAt = time.Now()
		s.log.Info().Str("pkg", req.Package).Msgf("registering app")
	} else {
		s.log.Info().Str("pkg", req.Package).Msgf("updating app")
	}

	ent.ID = appKey
	ent.Name = req.Name
	ent.Package = req.Package
	ent.Version = req.Version
	ent.APIURL = req.APIURL
	ent.WebURL = req.WebURL
	ent.RemoteEntry = req.RemoteEntryFile
	ent.Proxy = req.Proxy
	ent.Navigation = []*apptypes.Navigation{}
	ent.UpdatedAt = time.Now()
	ent.Adopted = true
	ent.Available = true
	ent.RemoteEntryRewriteRegEx = map[string]string{
		"/app/*/*": "/$2",
	}

	for _, navigation := range req.Navigation {
		n := &apptypes.Navigation{
			ID: hashutil.GetHash64([]byte(req.Package + ":" + navigation.Module.Path)),
		}

		apputil.MapNavInputToNavEntity(navigation, n)

		ent.Navigation = append(ent.Navigation, n)

		if navigation.Proxy {
			n.RemoteEntry = fmt.Sprintf("%s/app/%s/remoteEntry.js", "", ent.ID)
		} else {
			n.RemoteEntry = fmt.Sprintf("%s/%s", req.WebURL, req.RemoteEntryFile)
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

	st := rc.HSet(ctx, "apps", ent.ID, ent)
	if st.Err() != nil {
		s.log.Error().Err(st.Err()).Msgf("failed to cache application")
	}

	return &model.RegisterAppOutput{ID: ent.ID}, s.rebuildNavigation()
}

// KeepAlive monitors the healthiness of a remove application, after registering apps must send a keep alive
// this refreshes the ttl on the cache entry preventing the app from being marked as unavailable
func (s *service) KeepAlive(ctx context.Context, req *model.KeepAliveAppInput) (*model.KeepAliveAppOutput, error) {
	var (
		rc        = s.opts.RedisUseCase.Client()
		err       error
		appExists bool
		appKey    = s.appKeepAliveKey(req.Name)
	)

	er := rc.Exists(ctx, appKey)
	if err = er.Err(); err != nil {
		s.log.Error().Str("package", req.Pkg).Err(err).Msgf("failed to retrieve check for cached app entry")
		return nil, fmt.Errorf("failed to retrieve check for cached app entry: %w", err)
	}
	appExists = er.Val() > 0

	if appExists {
		cmd := rc.Expire(ctx, appKey, apptypes.KeepAliveTTL)

		if cmd.Err() != nil {
			s.log.Error().Str("package", req.Pkg).Err(err).Msgf("could not refresh cached app expiry")
		}

		rsp := &model.KeepAliveAppOutput{
			RegistrationRequired: false,
			Ok:                   true,
		}
		return rsp, nil
	}

	return &model.KeepAliveAppOutput{
		RegistrationRequired: true,
		Ok:                   false,
	}, nil
}

/************************************************************************/
/* SHELL CONFIGURATION
/************************************************************************/

// GetAppConfiguration fetches the shell app configuration from the cache, does not build the configuration, that instead
// happens any time an app is added, removed or updated
func (s *service) GetAppConfiguration(ctx context.Context, tenant string) (*model.ShellConfiguration, error) {
	var (
		rc            = s.opts.RedisUseCase.Client()
		configuration model.ShellConfiguration
	)

	cmd := rc.Get(ctx, "shell:configuration")
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	if err := cmd.Scan(&configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}

/************************************************************************/
/* HELPERS
/************************************************************************/

// rebuildNavigation rebuilds the navigation structure for the shell and updates the entry in the cache, this process
// is contention free because it is and should only be run on the leader in the cluster
//
//nolint:prealloc
func (s *service) rebuildNavigation() error {
	var (
		rc   = s.opts.RedisUseCase.Client()
		nc   = s.opts.NatsUseCase.Client()
		apps []*apptypes.App
	)

	s.log.Info().Msgf("rebuilding navigation cache due to change event")

	iter := rc.HGetAll(s.opts.Context, "apps")
	if iter.Err() != nil {
		return iter.Err()
	}

	for _, val := range iter.Val() {
		var app apptypes.App
		if err := json.Unmarshal([]byte(val), &app); err != nil {
			s.log.Error().Err(err).Msgf("faled to unmarshal: \n%s", string(debug.Stack()))
			return err
		}
		apps = append(apps, &app)
	}

	mappings := apputil.MapAppsToNavigation(apps)
	cmd := rc.Set(s.opts.Context, "shell:configuration", mappings, 0)

	if cmd.Err() == nil {
		if err := nc.Publish(apptypes.ShellConfigurationUpdatedSubject, []byte("{}")); err != nil {
			s.log.Warn().Err(err).Msgf("failed to publish configuration rebuilt event")
		}
	}

	return cmd.Err()
}

// appKeepAliveKey generates a cache key for a keep alive that requires constant refresh otherwise it will trigger
// the keyspace watcher and eventually mark the app as unavailable
func (s *service) appKeepAliveKey(name string) string {
	return "app:keepalive:" + name
}

// markApplicationAsUnavailable marks an app as unavailable, this is called when an app stops sending
// keep alive messages for a period of time and before the navigation is rebuilt
func (s *service) markApplicationAsUnavailable(id string) error {
	var (
		rc  = s.opts.RedisUseCase.Client()
		app apptypes.App
	)

	cmd := rc.HGet(s.opts.Context, "apps", id)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if err := cmd.Scan(&app); err != nil {
		return err
	}

	app.Available = false

	setCmd := rc.HSet(s.opts.Context, "apps", id, app)

	return setCmd.Err()
}

/************************************************************************/
/* FACTORY
/************************************************************************/

func NewService(opts *apptypes.APIGatewayOptions, log zerolog.Logger) apptypes.InternalService {
	return &service{
		log:  log,
		opts: opts,
		targetCache: imcache.NewSharded[string, *apptypes.ProxyTarget](apptypes.CacheShards, imcache.DefaultStringHasher64{},
			imcache.WithCleanerOption[string, *apptypes.ProxyTarget](apptypes.CacheCleanupFreq),
			imcache.WithEvictionCallbackOption[string, *apptypes.ProxyTarget](func(key string, val *apptypes.ProxyTarget, reason imcache.EvictionReason) {
				if reason == imcache.EvictionReasonExpired {
					log.Info().Str("uri", key).Msgf("proxy target evicted from cache")
				}
			}),
		),
	}
}
