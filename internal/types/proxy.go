package apptypes

import (
	"github.com/labstack/echo/v4"
	"net/url"
	"regexp"
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
