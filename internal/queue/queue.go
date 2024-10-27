package queue

import (
	"container/list"
	"errors"
	"time"

	"github.com/orderly-queue/orderly/internal/command"
	"github.com/orderly-queue/orderly/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ErrEmptyQueue = errors.New("queue is empty")
	ErrUnknown    = errors.New("unknown error")
)

type Queue struct {
	list *list.List
}

func New() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (q *Queue) Len() uint {
	len, _ := measure(string(command.Len), func() (uint, error) {
		return uint(q.list.Len()), nil
	})
	return len
}

func (q *Queue) Push(item string) {
	measure(string(command.Push), func() (struct{}, error) {
		q.list.PushBack(item)
		return struct{}{}, nil
	})
}

func (q *Queue) Pop() (string, error) {
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
	measure(string(command.Drain), func() (struct{}, error) {
		q.list.Init()
		return struct{}{}, nil
	})
}

func (q *Queue) Snapshot() []string {
	panic("not implemented")
}

func measure[T any](method string, f func() (T, error)) (T, error) {
	start := time.Now()

	out, err := f()

	dur := time.Since(start)
	metrics.CommandSeconds.With(prometheus.Labels{"method": method}).Observe(dur.Seconds())

	return out, err
}
