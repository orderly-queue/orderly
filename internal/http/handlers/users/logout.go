package users

import (
	"net/http"

	"github.com/henrywhitaker3/go-template/internal/app"
	"github.com/henrywhitaker3/go-template/internal/http/common"
	"github.com/henrywhitaker3/go-template/internal/http/middleware"
	"github.com/labstack/echo/v4"
)

type LogoutHandler struct {
	app *app.App
}

func NewLogout(app *app.App) *LogoutHandler {
	return &LogoutHandler{app: app}
}

func (l *LogoutHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := common.GetToken(c.Request())
		if err := l.app.Jwt.InvalidateUser(c.Request().Context(), token); err != nil {
			return common.Stack(err)
		}
		return c.NoContent(http.StatusAccepted)
	}
}

func (l *LogoutHandler) Method() string {
	return http.MethodPost
}

func (l *LogoutHandler) Path() string {
	return "/auth/logout"
}

func (l *LogoutHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Authenticated(),
	}
}
