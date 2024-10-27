package middleware

import (
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/orderly-queue/orderly/internal/app"
	"github.com/orderly-queue/orderly/internal/http/common"
	"github.com/orderly-queue/orderly/internal/tracing"
	"github.com/labstack/echo/v4"
)

func User(app *app.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := common.GetToken(c.Request())
			if token == "" {
				return next(c)
			}

			user, err := app.Jwt.VerifyUser(c.Request().Context(), token)
			if err == nil {
				c.SetRequest(c.Request().WithContext(common.SetUser(c.Request().Context(), user)))
			}

			if user != nil {
				if app.Config.Telemetry.Sentry.Enabled {
					if hub := sentryecho.GetHubFromContext(c); hub != nil {
						hub.Scope().SetUser(sentry.User{
							ID:    user.ID.UUID().String(),
							Email: user.Email,
							Name:  user.Name,
						})
					}
				}
				if app.Config.Telemetry.Tracing.Enabled {
					tracing.AddString(c.Request().Context(), "user_id", user.ID.UUID().String())
					tracing.AddString(c.Request().Context(), "request_id", common.RequestID(c))
				}
			}

			// TODO: add handling for cookies

			return next(c)
		}
	}
}
