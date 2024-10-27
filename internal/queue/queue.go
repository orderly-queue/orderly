package queue

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/orderly-queue/orderly/internal/command"
	"github.com/orderly-queue/orderly/internal/logger"
	"github.com/orderly-queue/orderly/internal/metrics"
	"github.com/orderly-queue/orderly/internal/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ErrEmptyQueue = errors.New("queue is empty")
	ErrUnknown    = errors.New("unknown error")
)

type Queue struct {
	list *list.List
	mu   *sync.RWMutex

	listenLock *sync.RWMutex
	listeners  map[uuid.UUID]chan struct{}
}

func New() *Queue {
	return &Queue{
		list:       list.New(),
		mu:         &sync.RWMutex{},
		listenLock: &sync.RWMutex{},
		listeners:  make(map[uuid.UUID]chan struct{}),
	}
}

func (q *Queue) Len() uint {
	q.mu.RLock()
	defer q.mu.RUnlock()
	len, _ := measure(string(command.Len), func() (uint, error) {
		return uint(q.list.Len()), nil
	})
	return len
}

func (q *Queue) Push(item string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	measure(string(command.Push), func() (struct{}, error) {
		q.list.PushBack(item)
		return struct{}{}, nil
	})
}

func (q *Queue) Pop() (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return measure(string(command.Pop), func() (string, error) {
		front := q.list.Front()
		if front == nil {
			return "", ErrEmptyQueue
		}
		item, ok := q.list.Remove(front).(string)
		if !ok {
			return "", ErrUnknown
		}
		return item, nil
	})
}

func (q *Queue) Drain() {
	q.mu.Lock()
	defer q.mu.Unlock()
	measure(string(command.Drain), func() (struct{}, error) {
		q.list.Init()
		return struct{}{}, nil
	})
}

func (q *Queue) Consume(ctx context.Context) (<-chan string, error) {
	id, _, err := q.listen()
	if err != nil {
		return nil, err
	}
	defer q.ignore(id)
	out := make(chan string, 100)

	go func() {
		metrics.Consumers.Inc()
		defer metrics.Consumers.Dec()
		tick := time.NewTicker(time.Millisecond)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				q.tryPop(ctx, out)
			}
		}
	}()

	return out, nil
}

func (q *Queue) tryPop(ctx context.Context, out chan<- string) {
	item, err := q.Pop()
	if err != nil {
		if !errors.Is(err, ErrEmptyQueue) {
			logger.Logger(ctx).Errorw("failed to pop from queue", "error", err)
		}
		return
	}
	out <- item
}

// Returns a channel that notifies when an item is pushed to the queue
func (q *Queue) listen() (uuid.UUID, <-chan struct{}, error) {
	id, err := uuid.New()
	if err != nil {
		return id, nil, err
	}

	ch := make(chan struct{}, 1)
	q.listenLock.Lock()
	q.listeners[id] = ch
	q.listenLock.Unlock()

	return id, ch, nil
}

func (q *Queue) ignore(id uuid.UUID) {
	q.listenLock.Lock()
	defer q.listenLock.Unlock()
	ch := q.listeners[id]
	delete(q.listeners, id)
	close(ch)
}

func (q *Queue) Snapshot() []string {
	out, _ := measure("snapshot", func() ([]string, error) {
		out := []string{}
		q.mu.Lock()
		defer q.mu.Unlock()
		for e := q.list.Front(); e != nil; e = e.Next() {
			o, ok := e.Value.(string)
			if ok {
				out = append(out, o)
			}
		}
		return out, nil
	})
	return out
}

// Empties the queue and loads the data into it form the input slice
func (q *Queue) Load(data []string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, d := range data {
		if d != "" {
			q.list.PushBack(d)
		}
	}
}

func measure[T any](method string, f func() (T, error)) (T, error) {
	start := time.Now()

	out, err := f()

	dur := time.Since(start)
	metrics.CommandSeconds.With(prometheus.Labels{"method": method}).Observe(dur.Seconds())

	return out, err
}
