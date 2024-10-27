package users

import (
	"net/http"

	"github.com/orderly-queue/orderly/internal/app"
	"github.com/orderly-queue/orderly/internal/http/common"
	"github.com/orderly-queue/orderly/internal/http/middleware"
	"github.com/labstack/echo/v4"
)

type MeHandler struct {
	app *app.App
}

func NewMe(app *app.App) *MeHandler {
	return &MeHandler{app: app}
}

func (m *MeHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := common.GetUser(c.Request().Context())
		return c.JSON(http.StatusOK, user)
	}
}

func (m *MeHandler) Method() string {
	return http.MethodGet
}

func (m *MeHandler) Path() string {
	return "/auth/me"
}

func (m *MeHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Authenticated(),
	}
}
