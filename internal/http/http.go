package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/go-template/internal/app"
	"github.com/henrywhitaker3/go-template/internal/http/handlers"
	"github.com/henrywhitaker3/go-template/internal/http/middleware"
	"github.com/henrywhitaker3/go-template/internal/logger"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

type Http struct {
	e   *echo.Echo
	app *app.App
}

func New(app *app.App) *Http {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	if app.Config.Telemetry.Tracing.Enabled {
		e.Use(middleware.Tracing())
	}
	e.Use(mw.RequestID())
	e.Use(middleware.Zap(app.Config.LogLevel.Level()))
	e.Use(mw.Recover())
	e.Use(middleware.Logger())

	h := &Http{
		e:   e,
		app: app,
	}

	h.Register(handlers.NewHello())

	return h
}

func (h *Http) Start(ctx context.Context) error {
	logger.Logger(ctx).Infow("starting http server", "port", h.app.Config.Http.Port)
	if err := h.e.Start(fmt.Sprintf(":%d", h.app.Config.Http.Port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (h *Http) Stop(ctx context.Context) error {
	logger.Logger(ctx).Info("stopping http server")
	return h.e.Shutdown(ctx)
}

func (h *Http) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.e.ServeHTTP(w, r)
}

type Handler interface {
	Handler() echo.HandlerFunc
	Method() string
	Path() string
	Middleware() []echo.MiddlewareFunc
}

func (h *Http) Register(handler Handler) {
	var reg func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

	switch handler.Method() {
	case http.MethodGet:
		reg = h.e.GET
	case http.MethodPost:
		reg = h.e.POST
	case http.MethodPatch:
		reg = h.e.PATCH
	case http.MethodDelete:
		reg = h.e.DELETE
	default:
		panic("invalid http method registered")
	}

	mw := handler.Middleware()
	if len(mw) == 0 {
		// Add a empty middleware so []... doesn't add a nil item
		mw = []echo.MiddlewareFunc{
			func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					return next(c)
				}
			},
		}
	}

	reg(handler.Path(), handler.Handler(), mw...)
}
