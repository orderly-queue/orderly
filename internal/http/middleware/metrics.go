package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/orderly-queue/orderly/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics(conf config.Telemetry, reg prometheus.Registerer) echo.MiddlewareFunc {
	return echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Subsystem: strings.ReplaceAll(conf.Tracing.ServiceName, "-", "_"),
		Skipper: func(c echo.Context) bool {
			return c.Request().Method == http.MethodOptions
		},
		Registerer: reg,
	})
}
