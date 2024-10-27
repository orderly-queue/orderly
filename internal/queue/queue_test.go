package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItPushesToTheQueue(t *testing.T) {
	queue := New()
	require.Equal(t, uint(0), queue.Len())
	queue.Push("bongo")
	require.Equal(t, uint(1), queue.Len())
}

func TestItPopsFromTheQueue(t *testing.T) {
	queue := New()
	require.Equal(t, uint(0), queue.Len())
	queue.Push("bongo")
	require.Equal(t, uint(1), queue.Len())
	item, err := queue.Pop()
	require.Nil(t, err)
	require.Equal(t, "bongo", item)
	require.Equal(t, uint(0), queue.Len())
}

func TestItDrainsTheQueue(t *testing.T) {
	queue := New()
	require.Equal(t, uint(0), queue.Len())
	queue.Push("bongo")
	require.Equal(t, uint(1), queue.Len())
	queue.Drain()
	require.Equal(t, uint(0), queue.Len())
}

func BenchmarkQueuePush(b *testing.B) {
	queue := New()

	item := "bongo"

	b.Run("push", func(b *testing.B) {
		for range b.N {
			queue.Push(item)
		}
	})
	b.Run("pop", func(b *testing.B) {
		for range b.N {
			queue.Push(item)
		}
		b.ResetTimer()
		b.ReportAllocs()
		for range b.N {
			_, err := queue.Pop()
			require.Nil(b, err)
		}
	})
}
