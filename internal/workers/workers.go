package workers

import (
	"context"
	"time"

	"github.com/henrywhitaker3/go-template/internal/logger"
	"github.com/henrywhitaker3/go-template/internal/metrics"
	"github.com/henrywhitaker3/go-template/internal/redis"
	"github.com/redis/rueidis"
)

type Worker interface {
	Name() string
	Run(ctx context.Context) error
	Interval() time.Duration
	Timeout() time.Duration
}

type Runner struct {
	redis rueidis.Client
	ctx   context.Context
}

func NewRunner(ctx context.Context, redis rueidis.Client) *Runner {
	return &Runner{
		redis: redis,
		ctx:   ctx,
	}
}

func (r *Runner) Register(w Worker) {
	tick := time.NewTicker(w.Interval())
	logger := logger.Logger(r.ctx)

	logger.Infow("registering worker", "name", w.Name(), "interval", w.Interval().String(), "timeout", w.Timeout().String())

	go func() {
		defer tick.Stop()
		for {
			select {
			case <-r.ctx.Done():
				logger.Infow("stopping worker", "name", w.Name())
				return
			case <-tick.C:
				ctx, cancel := context.WithTimeout(r.ctx, w.Timeout())

				if err := redis.Lock(ctx, r.redis, w.Name(), w.Timeout()); err != nil {
					logger.Debugw("skipping executing worker", "name", w.Name(), "error", err)
					// We didn't get the lock, so skip
					cancel()
					continue
				}

				metrics.WorkerExecutions.WithLabelValues(w.Name()).Inc()
				logger.Debugw("executing worker", "name", w.Name())

				if err := w.Run(ctx); err != nil {
					logger.Errorw("worker run failed", "name", w.Name(), "error", err)
					metrics.WorkerExecutionErrors.WithLabelValues(w.Name()).Inc()
				}

				// Release the lock, use the parent context so we still remove the
				// lock even if the child context has timed out
				redis.Unlock(r.ctx, r.redis, w.Name())
				cancel()
			}
		}
	}()
}
