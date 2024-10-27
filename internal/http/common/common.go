package common

import (
	"context"
	"net/http"
	"strings"

	"github.com/henrywhitaker3/ctxgen"
	"github.com/labstack/echo/v4"
)

var (
	ctxIdKey   = "request_id"
	traceIdKey = "trace_id"
)

func RequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func SetContextID(ctx context.Context, id string) context.Context {
	return ctxgen.WithValue(ctx, ctxIdKey, id)
}

func ContextID(ctx context.Context) string {
	return ctxgen.Value[string](ctx, ctxIdKey)
}

func SetTraceID(ctx context.Context, id string) context.Context {
	return ctxgen.WithValue(ctx, traceIdKey, id)
}

func TraceID(ctx context.Context) string {
	return ctxgen.Value[string](ctx, traceIdKey)
}

func SetRequest[T any](ctx context.Context, req T) context.Context {
	return ctxgen.WithValue(ctx, "request", req)
}

func GetRequest[T any](ctx context.Context) (T, bool) {
	return ctxgen.ValueOk[T](ctx, "request")
}

func GetToken(req *http.Request) string {
	header := req.Header.Get(echo.HeaderAuthorization)
	if header != "" {
		header = strings.Replace(header, "Bearer ", "", 1)
		return header
	}

	cookie, err := req.Cookie("auth")
	if err == nil {
		return cookie.Value
	}
	return ""
}
