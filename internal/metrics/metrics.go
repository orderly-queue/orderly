package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/henrywhitaker3/go-template/internal/logger"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	e *echo.Echo

	port int

	registry *prometheus.Registry
	reg      *sync.Once
}

func New(port int) *Metrics {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	m := &Metrics{
		e:        e,
		port:     port,
		registry: prometheus.NewRegistry(),
		reg:      &sync.Once{},
	}

	m.reg.Do(func() {})

	m.e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
		Gatherer: m.registry,
	}))

	return m
}

func (m *Metrics) Start(ctx context.Context) error {
	logger.Logger(ctx).Infow("starting metrics server", "port", m.port)
	if err := m.e.Start(fmt.Sprintf(":%d", m.port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (m *Metrics) Stop(ctx context.Context) error {
	logger.Logger(ctx).Info("stopping metrics server")
	return m.e.Shutdown(ctx)
}
