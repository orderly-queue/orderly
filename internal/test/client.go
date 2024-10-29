package test

import (
	"context"
	"testing"
	"time"

	"github.com/orderly-queue/orderly/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func Client(t *testing.T, timeout time.Duration) (context.Context, *sdk.Client, context.CancelFunc) {
	app, cancelApp := App(t, true)

	ctx, cancelCtx := context.WithTimeout(context.Background(), timeout)

	srv, cancelSrv := Server(app)

	client, err := sdk.NewClient(ctx, sdk.ClientConfig{
		Endpoint: srv,
	})
	require.Nil(t, err)

	return ctx, client, func() {
		require.Nil(t, client.Close())
		cancelSrv()
		cancelCtx()
		cancelApp()
	}
}
