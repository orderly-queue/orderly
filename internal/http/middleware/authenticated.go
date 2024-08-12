package middleware

import (
	"github.com/henrywhitaker3/go-template/internal/http/common"
	"github.com/labstack/echo/v4"
)

func Authenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, ok := common.GetUser(c.Request().Context())
			if !ok {
				return common.Stack(common.ErrUnauth)
			}
			return next(c)
		}
	}
}
