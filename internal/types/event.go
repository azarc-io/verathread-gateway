package types

type (
	AppAddedOrUpdatedEvent struct {
		APIEndpoint   string       `json:"apiUrl"`
		APIWsEndpoint string       `json:"apiWsUrl"`
		Package       string       `json:"package"`
		ProxyAPI      bool         `json:"proxyApi"`
		Name          string       `json:"name"`
		Version       string       `json:"version"`
		Available     bool         `json:"available"`
		Navigation    []*AppModule `json:"navigation"`
	}

	AppRemovedEvent struct {
		Package    string       `json:"package"`
		Name       string       `json:"name"`
		Version    string       `json:"version"`
		Navigation []*AppModule `json:"navigation"`
	}

	AppModule struct {
		ID                      string            `json:"id"`
		Proxy                   bool              `json:"proxy"`
		BaseURL                 string            `json:"baseUrl"`
		RemoteEntryRewriteRegEx map[string]string `json:"remoteEntryRewriteRegEx"`
		Slug                    string            `json:"slug"`
	}
)
