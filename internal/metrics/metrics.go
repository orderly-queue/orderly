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
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	CommandSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "orderly_command_seconds",
		Help: "The duration of commands processed by the server",
		Buckets: []float64{
			toSeconds(time.Nanosecond * 5), toSeconds(time.Nanosecond * 10), toSeconds(time.Nanosecond * 20),
			toSeconds(time.Nanosecond * 50), toSeconds(time.Nanosecond * 100), toSeconds(time.Nanosecond * 105),
			toSeconds(time.Nanosecond * 110), toSeconds(time.Nanosecond * 112), toSeconds(time.Nanosecond * 114),
			toSeconds(time.Nanosecond * 116), toSeconds(time.Nanosecond * 118), toSeconds(time.Nanosecond * 120),
			toSeconds(time.Nanosecond * 121), toSeconds(time.Nanosecond * 122), toSeconds(time.Nanosecond * 123),
			toSeconds(time.Nanosecond * 124), toSeconds(time.Nanosecond * 125), toSeconds(time.Nanosecond * 126),
			toSeconds(time.Nanosecond * 127), toSeconds(time.Nanosecond * 128), toSeconds(time.Nanosecond * 129),
			toSeconds(time.Nanosecond * 130), toSeconds(time.Nanosecond * 132), toSeconds(time.Nanosecond * 134),
			toSeconds(time.Nanosecond * 136), toSeconds(time.Nanosecond * 138), toSeconds(time.Nanosecond * 140),
			toSeconds(time.Nanosecond * 150), toSeconds(time.Nanosecond * 160), toSeconds(time.Nanosecond * 170),
			toSeconds(time.Nanosecond * 180), toSeconds(time.Nanosecond * 190), toSeconds(time.Nanosecond * 200),
			toSeconds(time.Nanosecond * 300), toSeconds(time.Nanosecond * 400), toSeconds(time.Nanosecond * 500),
		},
	}, []string{"method"})

	Consumers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "orderly_consumers",
		Help: "The current number of ocnnected consumers",
	})

	Size = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "orderly_queue_size",
		Help: "The size of the queue",
	})
	Pending = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "orderly_pending_notifications",
		Help: "The number of pending notifications for consumers",
	})
)

type Metrics struct {
	e *echo.Echo

	port int

	Registry prometheus.Registerer
	Gatherer prometheus.Gatherer
	reg      *sync.Once
}

func New(port int) *Metrics {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	reg := prometheus.NewRegistry()

	m := &Metrics{
		e:        e,
		port:     port,
		Registry: reg,
		Gatherer: reg,
		reg:      &sync.Once{},
	}

	m.reg.Do(func() {
		m.Registry.MustRegister(CommandSeconds)
		m.Registry.MustRegister(Consumers)
		m.Registry.MustRegister(Size)
		m.Registry.MustRegister(Pending)
		m.Registry.MustRegister(collectors.NewBuildInfoCollector())
		m.Registry.MustRegister(collectors.NewGoCollector())
		m.Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	})

	m.e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
		Gatherer: m.Gatherer,
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
	return dur.Seconds()
}
