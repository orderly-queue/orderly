package middleware

import (
	"github.com/henrywhitaker3/go-template/internal/app"
	"github.com/henrywhitaker3/go-template/internal/http/common"
	"github.com/labstack/echo/v4"
)

func Admin(app *app.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := common.GetUser(c.Request().Context())
			if !ok {
				return common.Stack(common.ErrUnauth)
			}
			user, err := app.Users.Get(c.Request().Context(), user.ID)
			if err != nil {
				return common.Stack(err)
			}
			if !user.Admin {
				return common.Stack(common.ErrForbidden)
			}

			return next(c)
		}
	}
}
