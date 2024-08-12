package jwt_test

import (
	"context"
	"testing"
	"time"

	"github.com/henrywhitaker3/go-template/internal/test"
	"github.com/henrywhitaker3/go-template/internal/users"
	"github.com/henrywhitaker3/go-template/internal/uuid"
	"github.com/stretchr/testify/require"
)

func TestItCreatesAUserJwt(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := &users.User{
		ID:        uuid.MustNew(),
		Name:      test.Word(),
		Email:     test.Email(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	token, err := app.Jwt.NewForUser(user, time.Second)
	require.Nil(t, err)

	valid, err := app.Jwt.VerifyUser(ctx, token)
	require.Nil(t, err)
	require.Equal(t, user.ID, valid.ID)
	require.Equal(t, user.Name, valid.Name)
	require.Equal(t, user.Email, valid.Email)

	// Test it fails validation after it has expired
	time.Sleep(time.Second * 2)

	_, err = app.Jwt.VerifyUser(ctx, token)
	require.NotNil(t, err)
}

func TestItFailsWhenTokenInvalidated(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := &users.User{
		ID:        uuid.MustNew(),
		Name:      test.Word(),
		Email:     test.Email(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	token, err := app.Jwt.NewForUser(user, time.Second*5)
	require.Nil(t, err)

	valid, err := app.Jwt.VerifyUser(ctx, token)
	require.Nil(t, err)
	require.Equal(t, user.ID, valid.ID)
	require.Equal(t, user.Name, valid.Name)
	require.Equal(t, user.Email, valid.Email)

	require.Nil(t, app.Jwt.InvalidateUser(ctx, token))

	_, err = app.Jwt.VerifyUser(ctx, token)
	require.NotNil(t, err)
}
