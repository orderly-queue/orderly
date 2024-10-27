package users_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/orderly-queue/orderly/internal/jwt"
	"github.com/orderly-queue/orderly/internal/test"
	"github.com/orderly-queue/orderly/internal/users"
	"github.com/orderly-queue/orderly/internal/uuid"
	"github.com/stretchr/testify/require"
)

func TestItGetsTheCurrentUser(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	user, _ := test.User(t, app)

	token, err := app.Jwt.NewForUser(user, time.Minute)
	require.Nil(t, err)

	type testCase struct {
		name  string
		user  *users.User
		token string
		code  int
	}

	tempJwt := jwt.New("bongo", app.Redis)

	randToken, err := tempJwt.NewForUser(&users.User{
		ID:    uuid.MustNew(),
		Email: test.Email(),
		Name:  test.Word(),
	}, time.Minute)
	require.Nil(t, err)

	tcs := []testCase{
		{
			name:  "gets the user for the token",
			user:  user,
			token: token,
			code:  http.StatusOK,
		},
		{
			name:  "returns a 401 when invalid token",
			user:  user,
			token: randToken,
			code:  http.StatusUnauthorized,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			rec := test.Get(app, "/auth/me", c.token)
			require.Equal(t, c.code, rec.Code)
			if c.code == http.StatusOK {
				var resp users.User
				require.Nil(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				require.Equal(t, c.user.ID, resp.ID)
				require.Equal(t, c.user.Name, resp.Name)
				require.Equal(t, c.user.Email, resp.Email)
			}
		})
	}
}
