package snapshotter_test

import (
	"context"
	"testing"
	"time"

	"github.com/orderly-queue/orderly/internal/test"
	"github.com/orderly-queue/orderly/internal/uuid"
	"github.com/stretchr/testify/require"
)

func TestItBacksupAQueue(t *testing.T) {
	app, cancel := test.App(t, true)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	id := uuid.MustNew().UUID().String()
	app.Queue.Push(id)

	require.Nil(t, app.Snapshotter.Snapshot(ctx))

	app.Queue.Drain()
	require.Equal(t, uint(0), app.Queue.Len())

	state, err := app.Snapshotter.Latest(ctx)
	require.Nil(t, err)
	require.Len(t, state, 1)

	app.Queue.Load(state)

	require.Equal(t, uint(1), app.Queue.Len())
	out, err := app.Queue.Pop()
	require.Nil(t, err)
	require.Equal(t, id, out)
}
