package users

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/go-template/internal/app"
	"github.com/henrywhitaker3/go-template/internal/http/common"
	"github.com/henrywhitaker3/go-template/internal/http/middleware"
	"github.com/labstack/echo/v4"
)

type RemoveAdminHandler struct {
	app *app.App
}

func NewRemoveAdmin(app *app.App) *RemoveAdminHandler {
	return &RemoveAdminHandler{app: app}
}

func (m *RemoveAdminHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req, ok := common.GetRequest[AdminRequest](c.Request().Context())
		if !ok {
			return common.ErrBadRequest
		}

		user, err := m.app.Users.Get(c.Request().Context(), req.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("%w: user not found", common.ErrValidation)
			}
			return common.Stack(err)
		}

		if err := m.app.Users.RemoveAdmin(c.Request().Context(), user); err != nil {
			return common.Stack(err)
		}

		return c.NoContent(http.StatusAccepted)
	}
}

func (m *RemoveAdminHandler) Method() string {
	return http.MethodDelete
}

func (m *RemoveAdminHandler) Path() string {
	return "/auth/admin"
}

func (m *RemoveAdminHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Bind[AdminRequest](),
		middleware.Authenticated(),
		middleware.Admin(m.app),
	}
}
