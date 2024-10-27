package middleware

import (
	"github.com/orderly-queue/orderly/internal/http/common"
	"github.com/orderly-queue/orderly/internal/logger"
	"github.com/orderly-queue/orderly/internal/tracing"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Zap(level zap.AtomicLevel) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = logger.Wrap(ctx, level)

			id := common.RequestID(c)
			if id != "" {
				ctx = common.SetContextID(ctx, id)
			}

			if traceId := tracing.TraceID(ctx); traceId != "" {
				ctx = common.SetTraceID(ctx, traceId)
			}

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
