package sdk_test

import (
	"context"
	"testing"
	"time"

	"github.com/orderly-queue/orderly/internal/test"
	"github.com/stretchr/testify/require"
)

func TestItConnects(t *testing.T) {
	_, _, cancel := test.Client(t, time.Second*5)
	defer cancel()
}

func TestItPushesToTheQueue(t *testing.T) {
	ctx, client, cancel := test.Client(t, time.Second*5)
	defer cancel()

	require.Nil(t, client.Push(ctx, test.Word()))
}

func TestItPopsFromTheQueue(t *testing.T) {
	ctx, client, cancel := test.Client(t, time.Second*5)
	defer cancel()

	item := test.Word()
	require.Nil(t, client.Push(ctx, item))

	time.Sleep(time.Millisecond)

	out, err := client.Pop(ctx)
	require.Nil(t, err)
	require.Equal(t, item, out)
}

func TestItConsumesFromTheQueue(t *testing.T) {
	ctx, client, cancel := test.Client(t, time.Second*5)
	defer cancel()

	counter := 0
	for range 5 {
		require.Nil(t, client.Push(ctx, test.Word()))
	}

	cons, err := client.Consume(ctx)
	require.Nil(t, err)
	go func() {
		for range cons {
			counter++
		}
	}()

	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("counter has not reached 5 before timeout")
	default:
		if counter == 5 {
			return
		}
	}
}
