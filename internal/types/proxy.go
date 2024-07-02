package apptypes

import (
	"net/url"
	"regexp"

	"github.com/labstack/echo/v4"
)

type (
	// ProxyTarget defines the upstream target.
	ProxyTarget struct {
		Name         string
		URL          *url.URL
		Meta         echo.Map
		RegexRewrite map[*regexp.Regexp]string
	}
)
