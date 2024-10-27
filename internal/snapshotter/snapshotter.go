package snapshotter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/orderly-queue/orderly/internal/config"
	"github.com/orderly-queue/orderly/internal/logger"
	"github.com/orderly-queue/orderly/internal/queue"
	"github.com/thanos-io/objstore"
)

type Snapshotter struct {
	queue  *queue.Queue
	bucket objstore.Bucket
	conf   config.Snapshot
}

func New(conf config.Snapshot, queue *queue.Queue, bucket objstore.Bucket) *Snapshotter {
	return &Snapshotter{
		conf:   conf,
		bucket: bucket,
		queue:  queue,
	}
}

func (s *Snapshotter) Work(ctx context.Context) error {
	logger := logger.Logger(ctx)
	logger.Infow("starting snapshotter", "schedule", s.conf.Schedule, "retention_days", s.conf.RetentionDays)

	sched, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	_, err = sched.NewJob(
		gocron.CronJob(s.conf.Schedule, false),
		gocron.NewTask(s.Snapshot, ctx),
	)
	if err != nil {
		return err
	}

	sched.Start()
	go func() {
		<-ctx.Done()
		sched.Shutdown()
		logger.Info("stopping snapshotter")
	}()
	return nil
}

func (s *Snapshotter) Snapshot(ctx context.Context) error {
	logger := logger.Logger(ctx)
	logger.Infow("snapshotting queue")
	data := s.queue.Snapshot()
	name := fmt.Sprintf("%s.state", time.Now().Format(time.RFC3339))
	by, err := json.Marshal(data)
	if err != nil {
		logger.Errorw("failed to marshall snapshot", "error", err)
		return err
	}
	if err := s.bucket.Upload(ctx, name, bytes.NewReader(by)); err != nil {
		logger.Errorw("failed to upload snapshot", "error", err)
		return err
	}
	return nil
}

func (s *Snapshotter) Latest(ctx context.Context) ([]string, error) {
	files := map[time.Time]string{}
	if err := s.bucket.Iter(ctx, "", func(name string) error {
		t, err := time.Parse(time.RFC3339, strings.ReplaceAll(name, ".state", ""))
		if err != nil {
			logger.Logger(ctx).Errorw("failed to parse snapshot name", "error", err)
			return nil
		}
		files[t] = name
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%w: failed to iterate over bucket", err)
	}

	if len(files) == 0 {
		return []string{}, nil
	}

	latest := time.Now().Add(-time.Hour * 24 * 365 * 10)
	for t := range files {
		if t.After(latest) {
			latest = t
		}
	}

	logger.Logger(ctx).Infow("loading snapshot", "snapshot", files[latest])

	state, err := s.bucket.Get(ctx, files[latest])
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get file", err)
	}
	read, err := io.ReadAll(state)
	if err != nil {
		return nil, err
	}

	out := []string{}
	if err := json.Unmarshal(read, &out); err != nil {
		return nil, err
	}
	return out, nil
}
