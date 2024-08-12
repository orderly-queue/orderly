package handlers

import (
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/go-template/internal/http/common"
	"github.com/henrywhitaker3/go-template/internal/http/middleware"
	"github.com/labstack/echo/v4"
)

type HelloRequest struct {
	Name string `param:"name"`
}

func (h HelloRequest) Validate() error {
	return nil
}

type Hello struct{}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req, ok := common.GetRequest[HelloRequest](c.Request().Context())
		if !ok {
			return common.ErrBadRequest
		}
		return c.String(http.StatusOK, fmt.Sprintf("Hello %s!", req.Name))
	}
}

func (h *Hello) Method() string {
	return http.MethodGet
}

func (h *Hello) Path() string {
	return "/:name"
}

func (h *Hello) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Bind[HelloRequest](),
	}
}
