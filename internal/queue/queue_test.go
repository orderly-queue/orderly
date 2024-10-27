package queue

import (
	"testing"

	"github.com/orderly-queue/orderly/internal/uuid"
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

func TestItSnapshots(t *testing.T) {
	queue := New()

	items := []string{}
	for range 10 {
		id := uuid.MustNew()
		items = append(items, id.UUID().String())
		queue.Push(id.UUID().String())
	}

	snap := queue.Snapshot()

	require.Len(t, snap, 10)
	require.Equal(t, items, snap)
}

func TestItSnapshotsEmptyList(t *testing.T) {
	queue := New()
	snap := queue.Snapshot()
	require.Len(t, snap, 0)
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
