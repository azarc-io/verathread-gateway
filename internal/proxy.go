package internal

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/azarc-io/verathread-gateway/internal/util"
	"github.com/erni27/imcache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	apptypes "github.com/azarc-io/verathread-gateway/internal/types"

	"github.com/labstack/echo/v4"
)

// StatusCodeContextCanceled is a custom HTTP status code for situations
// where a client unexpectedly closed the connection to the server.
// As there is no standard error code for "client closed connection", but
// various well-known HTTP clients and server implement this HTTP code we use
// 499 too instead of the more problematic 5xx, which does not allow to detect this situation
const StatusCodeContextCanceled = 499

// proxyCacheTimeoutDuration how long to store cached proxies for
var proxyCacheTimeoutDuration = time.Minute * 2

type (
	proxy struct {
		httpProxyCache *imcache.Sharded[string, *httputil.ReverseProxy]
		log            zerolog.Logger
		filesToScan    []string
	}
)

// responseModifier modifies proxied responses, unzipping them if required and scanning for tokens
func (p *proxy) responseModifier(response *http.Response, c echo.Context) error {
	var (
		encoding = response.Header.Get(echo.HeaderContentEncoding)
		body     []byte
		err      error
		reader   io.Reader
	)

	if strings.Contains(encoding, "gzip") || strings.Contains(encoding, "deflate") {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(reader)
		if err != nil {
			return err
		}
	} else {
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}
	}

	// will replace and cache tokens in served files
	if apputil.ShouldReplace(c.Request().URL.Path, p.filesToScan) {
		body = apputil.ReplaceTokens(c, body)
	}

	response.Body = io.NopCloser(bytes.NewBuffer(body))
	response.Header.Set("Cache-Control", "max-age=31536000")

	return nil
}

// errorHandler handles proxy errors
func (p *proxy) errorHandler(resp http.ResponseWriter, req *http.Request, err error, tgt *apptypes.ProxyTarget, c echo.Context) {
	target := c.Get(apptypes.TargetURLKey).(*url.URL)
	desc := target.String()
	if tgt.Name != "" {
		desc = fmt.Sprintf("%s(%s)", tgt.Name, desc)
	}
	// If the client canceled the request (usually by closing the connection), we can report a
	// client error (4xx) instead of a server error (5xx) to correctly identify the situation.
	// The Go standard library (as of late 2020) wraps the exported, standard
	// context.Canceled error with unexported garbage value requiring a substring check, see
	// https://github.com/golang/go/blob/6965b01ea248cabb70c3749fd218b36089a21efb/src/net/net.go#L416-L430
	if errors.Is(err, context.Canceled) || strings.Contains(err.Error(), "operation was canceled") {
		httpError := echo.NewHTTPError(StatusCodeContextCanceled, fmt.Sprintf("client closed connection: %v", err))
		httpError.Internal = err
		c.Set("_error", httpError)
	} else {
		httpError := echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("remote %s unreachable, could not forward: %v", desc, err))
		httpError.Internal = err
		c.Set("_error", httpError)
	}
}

// proxyRaw proxies raw TCP connection, capable of handling web sockets and SSE
func (p *proxy) proxyRaw(t *apptypes.ProxyTarget, c echo.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			target  = c.Get(apptypes.TargetURLKey).(*url.URL)
			in, out net.Conn
			err     error
		)

		in, _, err = c.Response().Hijack()
		if err != nil {
			c.Set("_error", fmt.Errorf("proxy raw, hijack error=%w, url=%s", err, target.String()))
			return
		}
		defer func(in net.Conn) {
			err = in.Close()
			if err != nil {
				p.log.Warn().Err(err).Msgf("error while closing inbound proxy connection")
			}
		}(in)

		out, err = net.Dial("tcp", target.Host)
		if err != nil {
			log.Warn().Err(err).Str("url", target.String()).Msgf("proxy raw, dial error")
			c.Set("_error", echo.NewHTTPError(http.StatusBadGateway,
				fmt.Sprintf("proxy raw, dial error=%v, url=%s", err, target.String())))
			return
		}
		defer func(out net.Conn) {
			err = out.Close()
			if err != nil {
				p.log.Warn().Err(err).Msgf("error while closing outbound proxy connection")
			}
		}(out)

		// Write header
		err = r.Write(out)
		if err != nil {
			c.Set("_error", echo.NewHTTPError(http.StatusBadGateway,
				fmt.Sprintf("proxy raw, request header copy error=%v, url=%s", err, target.String())))
			return
		}

		//nolint:mnd
		errCh := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err = io.Copy(dst, src)
			errCh <- err
		}

		go cp(out, in)
		go cp(in, out)
		err = <-errCh
		if err != nil && err != io.EOF {
			p.log.Warn().Err(err).Str("url", target.String()).Msgf("proxy raw, copy body error")
			c.Set("_error", fmt.Errorf("proxy raw, copy body error=%w, url=%s", err, target.String()))
		}
	})
}

// proxyHTTP proxies http requests, uses caching to improve performance and reduce memory allocations
func (p *proxy) proxyHTTP(tgt *apptypes.ProxyTarget, c echo.Context) http.Handler {
	target := c.Get(apptypes.TargetURLKey).(*url.URL)
	proxy, exists := p.httpProxyCache.Get(target.String())

	if !exists {
		log.Info().Str("target", target.String()).Msgf("creating new proxy for target service")
		proxy = httputil.NewSingleHostReverseProxy(target)
		p.httpProxyCache.Set(target.String(), proxy, imcache.WithExpiration(proxyCacheTimeoutDuration))
	}

	proxy.ModifyResponse = func(response *http.Response) error {
		return p.responseModifier(response, c)
	}

	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		p.errorHandler(resp, req, err, tgt, c)
	}

	return proxy
}

// rewriteURL applies any url re-write expressions
func (p *proxy) rewriteURL(rewriteRegex map[*regexp.Regexp]string, req *http.Request) error {
	if len(rewriteRegex) == 0 {
		return nil
	}

	// Depending on how HTTP request is sent RequestURI could contain Scheme://Host/path or be just /path.
	// We only want to use path part for rewriting and therefore trim prefix if it exists
	rawURI := req.RequestURI
	if rawURI != "" && rawURI[0] != '/' {
		prefix := ""
		if req.URL.Scheme != "" {
			prefix = req.URL.Scheme + "://"
		}
		if req.URL.Host != "" {
			prefix += req.URL.Host // host or host:port
		}
		if prefix != "" {
			rawURI = strings.TrimPrefix(rawURI, prefix)
		}
	}

	for k, v := range rewriteRegex {
		if replacer := p.captureTokens(k, rawURI); replacer != nil {
			url, err := req.URL.Parse(replacer.Replace(v))
			if err != nil {
				return err
			}
			req.URL = url

			return nil // rewrite only once
		}
	}
	return nil
}

// captureTokens captures url tokens for re-writing
func (p *proxy) captureTokens(pattern *regexp.Regexp, input string) *strings.Replacer {
	groups := pattern.FindAllStringSubmatch(input, -1)
	if groups == nil {
		return nil
	}

	values := groups[0][1:]
	//nolint:mnd
	replace := make([]string, 2*len(values))
	for i, v := range values {
		//nolint:mnd
		j := 2 * i
		replace[j] = "$" + strconv.Itoa(i+1)
		replace[j+1] = v
	}
	return strings.NewReplacer(replace...)
}
