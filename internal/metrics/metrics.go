package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/orderly-queue/orderly/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CommandSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "command_seconds",
		Help: "The duration of commands processed by the server",
		Buckets: []float64{
			toSeconds(time.Microsecond * 5), toSeconds(time.Microsecond * 10), toSeconds(time.Microsecond * 20),
			toSeconds(time.Microsecond * 50), toSeconds(time.Microsecond * 100), toSeconds(time.Microsecond * 200),
			toSeconds(time.Microsecond * 500), toSeconds(time.Millisecond),
		},
	}, []string{"method"})
)

type Metrics struct {
	e *echo.Echo

	port int

	Registry *prometheus.Registry
	reg      *sync.Once
}

func New(port int) *Metrics {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	m := &Metrics{
		e:        e,
		port:     port,
		Registry: prometheus.NewRegistry(),
		reg:      &sync.Once{},
	}

	m.reg.Do(func() {
		m.Registry.MustRegister(CommandSeconds)
	})

	m.e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
		Gatherer: m.Registry,
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

func toSeconds(dur time.Duration) float64 {
	return float64(dur) / float64(time.Second)
}
