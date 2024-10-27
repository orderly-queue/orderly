package users

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/orderly-queue/orderly/internal/app"
	"github.com/orderly-queue/orderly/internal/http/common"
	"github.com/orderly-queue/orderly/internal/http/middleware"
	"github.com/orderly-queue/orderly/internal/uuid"
	"github.com/labstack/echo/v4"
)

type AdminRequest struct {
	ID uuid.UUID `json:"id"`
}

func (a AdminRequest) Validate() error {
	return nil
}

type MakeAdminHandler struct {
	app *app.App
}

func NewMakeAdmin(app *app.App) *MakeAdminHandler {
	return &MakeAdminHandler{app: app}
}

func (m *MakeAdminHandler) Handler() echo.HandlerFunc {
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

		if err := m.app.Users.MakeAdmin(c.Request().Context(), user); err != nil {
			return common.Stack(err)
		}

		return c.NoContent(http.StatusAccepted)
	}
}

func (m *MakeAdminHandler) Method() string {
	return http.MethodPost
}

func (m *MakeAdminHandler) Path() string {
	return "/auth/admin"
}

func (m *MakeAdminHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Authenticated(),
		middleware.Admin(m.app),
		middleware.Bind[AdminRequest](),
	}
}
