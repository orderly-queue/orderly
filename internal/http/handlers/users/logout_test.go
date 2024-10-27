package users_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/orderly-queue/orderly/internal/test"
	"github.com/stretchr/testify/require"
)

func TestItLogsOutAUser(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	user, _ := test.User(t, app)

	token, err := app.Jwt.NewForUser(user, time.Minute)
	require.Nil(t, err)

	rec := test.Post(t, app, "/auth/logout", nil, token)

	require.Equal(t, http.StatusAccepted, rec.Code)

	rec = test.Get(app, "/auth/me", token)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
