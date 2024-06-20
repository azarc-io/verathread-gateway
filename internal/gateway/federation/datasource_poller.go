package federation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/wundergraph/graphql-go-tools/execution/engine"
)

type ServiceConfig struct {
	Name     string
	URL      string
	WS       string
	Fallback func(*ServiceConfig) (string, error)
}

type DatasourcePollerConfig struct {
	Services        []*ServiceConfig
	PollingInterval time.Duration
}

func NewDatasourcePoller(httpClient *http.Client, logger zerolog.Logger, config DatasourcePollerConfig) *DatasourcePollerPoller {
	return &DatasourcePollerPoller{
		httpClient: httpClient,
		config:     config,
		sdlMap:     make(map[string]string),
		log:        logger,
	}
}

type DatasourcePollerPoller struct {
	httpClient                *http.Client
	config                    DatasourcePollerConfig
	sdlMap                    map[string]string
	updateDatasourceObservers []DataSourceObserver
	log                       zerolog.Logger
	sync.Mutex
}

func (d *DatasourcePollerPoller) Register(updateDatasourceObserver DataSourceObserver) {
	d.updateDatasourceObservers = append(d.updateDatasourceObservers, updateDatasourceObserver)
}

func (d *DatasourcePollerPoller) Run(ctx context.Context) {
	d.updateSDLs(ctx)

	if d.config.PollingInterval == 0 {
		<-ctx.Done()
		return
	}

	ticker := time.NewTicker(d.config.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.updateSDLs(ctx)
		}
	}
}

func (d *DatasourcePollerPoller) updateSDLs(ctx context.Context) {
	d.sdlMap = make(map[string]string)

	d.Lock()
	defer d.Unlock()

	var wg sync.WaitGroup
	resultCh := make(chan struct {
		name string
		sdl  string
	})

	for _, serviceConf := range d.config.Services {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sdl, err := d.fetchServiceSDL(ctx, serviceConf.URL)
			if err != nil {
				log.Println("Failed to get sdl.", err)

				if serviceConf.Fallback == nil {
					return
				} else {
					sdl, err = serviceConf.Fallback(serviceConf)
					if err != nil {
						log.Println("Failed to get sdl with fallback.", err)
						return
					}
				}
			}

			select {
			case <-ctx.Done():
			case resultCh <- struct {
				name string
				sdl  string
			}{name: serviceConf.Name, sdl: sdl}:
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		d.sdlMap[result.name] = result.sdl
	}

	d.updateObservers()
}

func (d *DatasourcePollerPoller) updateObservers() {
	subgraphsConfig := d.createSubgraphsConfig()

	for i := range d.updateDatasourceObservers {
		d.updateDatasourceObservers[i].UpdateDataSources(subgraphsConfig)
	}
}

func (d *DatasourcePollerPoller) createSubgraphsConfig() []engine.SubgraphConfiguration {
	subgraphConfigs := make([]engine.SubgraphConfiguration, 0, len(d.config.Services))

	for _, serviceConfig := range d.config.Services {
		sdl, exists := d.sdlMap[serviceConfig.Name]
		if !exists {
			continue
		}

		subgraphConfig := engine.SubgraphConfiguration{
			Name:            serviceConfig.Name,
			URL:             serviceConfig.URL,
			SubscriptionUrl: serviceConfig.WS,
			SDL:             sdl,
		}

		subgraphConfigs = append(subgraphConfigs, subgraphConfig)
	}

	return subgraphConfigs
}

func (d *DatasourcePollerPoller) fetchServiceSDL(ctx context.Context, serviceURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serviceURL, bytes.NewReader([]byte(ServiceDefinitionQuery)))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	var result struct {
		Data struct {
			Service struct {
				SDL string `json:"sdl"`
			} `json:"_service"`
		} `json:"data"`
		Errors GQLErrors `json:"errors,omitempty"`
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read bytes: %w", err)
	}

	if err = json.NewDecoder(bytes.NewReader(bs)).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if result.Errors != nil {
		return "", fmt.Errorf("response error: %w", result.Errors)
	}

	return result.Data.Service.SDL, nil
}

func (d *DatasourcePollerPoller) addService(name string, url, wsURL string) {
	var exists bool

	for _, service := range d.config.Services {
		if service.Name == name {
			exists = true
			break
		}
	}

	if !exists {
		d.config.Services = append(d.config.Services, &ServiceConfig{
			Name:     name,
			URL:      url,
			WS:       wsURL,
			Fallback: nil,
		})

		sdl, err := d.fetchServiceSDL(context.Background(), url)
		if err != nil {
			d.log.Error().Err(err).
				Msgf("failed to fetch service sdl from <%s> skipping federation for this service", url)
		} else {
			d.sdlMap[name] = sdl
			d.updateObservers()
		}
	}
}

func (d *DatasourcePollerPoller) removeService(pkg string) {
	var exists bool

	for i, service := range d.config.Services {
		if service.Name == pkg {
			exists = true
			d.config.Services = append(d.config.Services[:i], d.config.Services[i+1:]...)
			break
		}
	}

	if exists {
		d.updateSDLs(context.Background())
	}
}
