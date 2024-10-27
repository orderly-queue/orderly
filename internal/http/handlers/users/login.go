package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/orderly-queue/orderly/internal/app"
	"github.com/orderly-queue/orderly/internal/http/common"
	"github.com/orderly-queue/orderly/internal/http/middleware"
	"github.com/orderly-queue/orderly/internal/metrics"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginRequest) Validate() error {
	if l.Email == "" {
		return fmt.Errorf("%w email", common.ErrRequiredField)
	}
	if l.Password == "" {
		return fmt.Errorf("%w password", common.ErrRequiredField)
	}
	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginHandler struct {
	app *app.App
}

func NewLogin(app *app.App) *LoginHandler {
	return &LoginHandler{app: app}
}

func (l *LoginHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req, ok := common.GetRequest[LoginRequest](c.Request().Context())
		if !ok {
			return common.ErrBadRequest
		}

		user, err := l.app.Users.Login(c.Request().Context(), req.Email, req.Password)
		if err != nil {
			metrics.Logins.WithLabelValues("false").Inc()
			return common.ErrUnauth
		}

		token, err := l.app.Jwt.NewForUser(user, time.Hour)
		if err != nil {
			return common.Stack(err)
		}

		metrics.Logins.WithLabelValues("true").Inc()

		return c.JSON(http.StatusOK, LoginResponse{
			Token: token,
		})
	}
}

func (l *LoginHandler) Method() string {
	return http.MethodPost
}

func (l *LoginHandler) Path() string {
	return "/auth/login"
}

func (l *LoginHandler) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Bind[LoginRequest](),
	}
}
