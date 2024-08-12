package users_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/henrywhitaker3/go-template/internal/http/handlers/users"
	"github.com/henrywhitaker3/go-template/internal/test"
	"github.com/stretchr/testify/require"
)

func TestItLogsInAUser(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, password := test.User(t, app)

	rec := test.Post(t, app, "/auth/login", users.LoginRequest{
		Email:    user.Email,
		Password: password,
	}, "")

	require.Equal(t, http.StatusOK, rec.Code)

	resp := &users.LoginResponse{}
	require.Nil(t, json.Unmarshal(rec.Body.Bytes(), resp))

	tuser, err := app.Jwt.VerifyUser(ctx, resp.Token)
	require.Nil(t, err)
	require.Equal(t, user.ID, tuser.ID)
}

func TestItErrorsWhenIncorrectPassword(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	user, _ := test.User(t, app)

	rec := test.Post(t, app, "/auth/login", users.LoginRequest{
		Email:    user.Email,
		Password: test.Sentence(5),
	}, "")

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestItErrorsWhenIncorrectEmail(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	rec := test.Post(t, app, "/auth/login", users.LoginRequest{
		Email:    test.Email(),
		Password: test.Sentence(5),
	}, "")

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
