package probes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/henrywhitaker3/go-template/internal/logger"
	"github.com/labstack/echo/v4"
)

type Probes struct {
	mu *sync.RWMutex

	e    *echo.Echo
	port int

	ready   bool
	healthy bool
}

func New(port int) *Probes {
	p := &Probes{
		port: port,
		mu:   &sync.RWMutex{},
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/readyz", p.readyHandler())
	e.GET("/healthz", p.healthyHandler())

	p.e = e

	return p
}

func (p *Probes) Healthy() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.healthy = true
}

func (p *Probes) Unhealthy() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.healthy = false
}

func (p *Probes) Ready() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ready = true
}

func (p *Probes) Unready() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ready = false
}

func (p *Probes) Start(ctx context.Context) error {
	logger.Logger(ctx).Infow("starting probes server", "port", p.port)
	if err := p.e.Start(fmt.Sprintf(":%d", p.port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (p *Probes) Stop(ctx context.Context) error {
	logger.Logger(ctx).Info("stopping probes server")
	return p.e.Shutdown(ctx)
}

func (p *Probes) readyHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		p.mu.RLock()
		defer p.mu.RUnlock()
		if p.ready {
			return c.String(http.StatusOK, "READY")
		}
		return c.String(http.StatusServiceUnavailable, "NOT READY")
	}
}

func (p *Probes) healthyHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		p.mu.RLock()
		defer p.mu.RUnlock()
		if p.healthy {
			return c.String(http.StatusOK, "HEALTHY")
		}
		return c.String(http.StatusServiceUnavailable, "NOT HEALTHY")
	}
}
