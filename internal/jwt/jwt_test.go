package jwt_test

import (
	"context"
	"testing"
	"time"

	"github.com/orderly-queue/orderly/internal/test"
	"github.com/stretchr/testify/require"
)

func TestItCreatesAJwt(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	token, err := app.Jwt.New(time.Second)
	require.Nil(t, err)

	err = app.Jwt.Verify(ctx, token)
	require.Nil(t, err)

	// Test it fails validation after it has expired
	time.Sleep(time.Second * 2)

	err = app.Jwt.Verify(ctx, token)
	require.NotNil(t, err)
}
