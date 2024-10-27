package middleware

import (
	"net/http"

	"github.com/orderly-queue/orderly/internal/http/common"
	"github.com/labstack/echo/v4"
)

type request interface {
	Validate() error
}

func Bind[T request]() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var req T
			if err := c.Bind(&req); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest)
			}
			if err := req.Validate(); err != nil {
				return err
			}
			ctx := common.SetRequest(c.Request().Context(), req)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
