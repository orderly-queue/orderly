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
	"github.com/orderly-queue/orderly/internal/logger"
	"github.com/orderly-queue/orderly/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thanos-io/objstore"
)

type store interface {
	Snapshot() []string
}

type Snapshotter struct {
	queue  store
	bucket objstore.Bucket
	conf   config.Snapshot

	age  *snapshotAge
	size *snapshotSize
}

func New(conf config.Snapshot, queue store, bucket objstore.Bucket, reg prometheus.Registerer) *Snapshotter {
	age := newSnapshotAge()
	size := newSnapshotSize()
	reg.MustRegister(age)
	reg.MustRegister(size)
	return &Snapshotter{
		conf:   conf,
		bucket: bucket,
		queue:  queue,
		age:    age,
		size:   size,
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
	_, err = sched.NewJob(
		gocron.CronJob("* * * * *", false),
		gocron.NewTask(s.report, ctx),
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

func (s *Snapshotter) report(ctx context.Context) error {
	logger := logger.Logger(ctx)
	logger.Debug("collecting snapshot metrics")

	snapshots, err := s.collect(ctx)
	if err != nil {
		logger.Errorw("failed to collect snapshots", "error", err)
		return err
	}

	s.age.record(snapshots)
	s.size.record(snapshots)

	return nil
}

func (s *Snapshotter) Latest(ctx context.Context) (*Snapshot, error) {
	files := map[time.Time]Snapshot{}
	snapshots, err := s.collect(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range snapshots {
		files[s.Time] = s
	}

	if len(files) == 0 {
		return nil, nil
	}

	latest := time.Now().Add(-time.Hour * 24 * 365 * 10)
	for t := range files {
		if t.After(latest) {
			latest = t
		}
	}
	l := files[latest]

	return &l, nil
}

func (s *Snapshotter) Open(ctx context.Context, snapshot Snapshot) ([]string, error) {
	logger := logger.Logger(ctx)
	logger.Infow("opening snapshot", "name", snapshot.Name)

	raw, err := s.bucket.Get(ctx, snapshot.Name)
	if err != nil {
		return nil, err
	}
	by, err := io.ReadAll(raw)
	if err != nil {
		return nil, err
	}
	out := []string{}
	if err := json.Unmarshal(by, &out); err != nil {
		return nil, err
	}
	return out, nil
}

type Snapshot struct {
	Time time.Time
	Name string
	Size int64
}

func (s Snapshot) Age() time.Duration {
	return time.Since(s.Time)
}

func (s *Snapshotter) collect(ctx context.Context) ([]Snapshot, error) {
	logger := logger.Logger(ctx)
	out := []Snapshot{}
	if err := s.bucket.Iter(ctx, "", func(name string) error {
		info, err := s.bucket.Attributes(ctx, name)
		if err != nil {
			logger.Errorw("failed to stat snapshot", "name", name, "error", err)
			return nil
		}
		t, err := time.Parse(time.RFC3339, strings.ReplaceAll(name, ".state", ""))
		if err != nil {
			logger.Errorw("failed to parse snapshot name", "error", err)
			return nil
		}

		out = append(out, Snapshot{
			Time: t,
			Name: name,
			Size: info.Size,
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}
