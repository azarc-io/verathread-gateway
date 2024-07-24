package middleware

import (
	"github.com/labstack/echo/v4"
)

// ACAOHeaderOverwriteMiddleware header entry de-duplication
// TODO make the keys configurable through the yaml config file
func ACAOHeaderOverwriteMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Before(func() {
			h := ctx.Response().Header()
			h.Set("X-Content-Type-Options", h.Values("X-Content-Type-Options")[0])
			h.Set("X-Dns-Prefetch-Control", h.Values("X-Dns-Prefetch-Control")[0])
			h.Set("X-Download-Options", h.Values("X-Download-Options")[0])
			h.Set("X-Frame-Options", h.Values("X-Frame-Options")[0])
			h.Set("X-Request-Id", h.Values("X-Request-Id")[0])
			h.Set("X-Xss-Protection", h.Values("X-Xss-Protection")[0])
			h.Set("Vary", h.Values("Vary")[0])
		})
		return next(ctx)
	}
}
