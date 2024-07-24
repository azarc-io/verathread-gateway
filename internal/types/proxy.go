package types

import (
	"encoding/json"
	"net/url"
	"regexp"

	"github.com/labstack/echo/v4"
)

type (
	// ProxyTarget defines the upstream target.
	ProxyTarget struct {
		ID           string
		Name         string
		WebURL       *url.URL
		APIURL       *url.URL
		Meta         echo.Map
		RegexRewrite map[*regexp.Regexp]string
	}
)

func (a ProxyTarget) MarshalBinary() (data []byte, err error) {
	return json.Marshal(a)
}

func (a *ProxyTarget) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}
