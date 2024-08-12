package middleware

import (
	"net/http"
	"slices"

	"github.com/henrywhitaker3/go-template/internal/config"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

var (
	skipPaths = []string{"/metrics", "/"}
)

func Tracing(conf config.Tracing) echo.MiddlewareFunc {
	return otelecho.Middleware(
		conf.ServiceName,
		otelecho.WithSkipper(func(c echo.Context) bool {
			if c.Request().Method == http.MethodOptions {
				return true
			}
			if slices.Contains(skipPaths, c.Path()) {
				return true
			}
			return false
		}),
	)
}
