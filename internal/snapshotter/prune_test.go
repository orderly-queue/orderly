package snapshotter

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/orderly-queue/orderly/internal/queue"
	"github.com/orderly-queue/orderly/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"github.com/thanos-io/objstore/providers/filesystem"
)

func TestItPrunesSnapshots(t *testing.T) {
	bucket, err := filesystem.NewBucket(t.TempDir())
	require.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	snap := New(config.Snapshot{
		Enabled:       true,
		Schedule:      "* * * *",
		RetentionDays: 1,
	}, queue.New(), bucket, prometheus.NewRegistry())

	twoDays := snap.name(time.Now().Add(-(time.Hour * 48)))
	hour := snap.name(time.Now().Add(-time.Hour))
	require.Nil(t, bucket.Upload(ctx, twoDays, bytes.NewReader([]byte("hello"))))
	require.Nil(t, bucket.Upload(ctx, hour, bytes.NewReader([]byte("goodbye"))))

	snaps, err := snap.collect(ctx)
	require.Nil(t, err)
	require.Len(t, snaps, 2)

	require.Nil(t, snap.prune(ctx))

	snaps, err = snap.collect(ctx)
	require.Nil(t, err)
	require.Len(t, snaps, 1)

	latest, err := snap.Latest(ctx)
	require.Nil(t, err)
	require.NotNil(t, latest)
	require.Equal(t, hour, latest.Name)
}

func TestItPrunesSnapshotsWithNamePrefix(t *testing.T) {
	bucket, err := filesystem.NewBucket(t.TempDir())
	require.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	snap := New(config.Snapshot{
		Enabled:       true,
		Schedule:      "* * * *",
		RetentionDays: 1,
		NamePrefix:    "bongo",
	}, queue.New(), bucket, prometheus.NewRegistry())

	twoDays := snap.name(time.Now().Add(-(time.Hour * 48)))
	hour := snap.name(time.Now().Add(-time.Hour))
	require.Nil(t, bucket.Upload(ctx, twoDays, bytes.NewReader([]byte("hello"))))
	require.Nil(t, bucket.Upload(ctx, hour, bytes.NewReader([]byte("goodbye"))))

	snaps, err := snap.collect(ctx)
	require.Nil(t, err)
	require.Len(t, snaps, 2)

	require.Nil(t, snap.prune(ctx))

	snaps, err = snap.collect(ctx)
	require.Nil(t, err)
	require.Len(t, snaps, 1)

	latest, err := snap.Latest(ctx)
	require.Nil(t, err)
	require.NotNil(t, latest)
	require.Equal(t, hour, latest.Name)
}
