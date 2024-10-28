package snapshotter

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type snapshotAge struct {
	mu        *sync.Mutex
	Desc      *prometheus.Desc
	snapshots []Snapshot
}

func newSnapshotAge() *snapshotAge {
	return &snapshotAge{
		mu: &sync.Mutex{},
		Desc: prometheus.NewDesc(
			"orderly_snapshots_age",
			"The age of the snapshots",
			[]string{"name"},
			nil,
		),
		snapshots: []Snapshot{},
	}
}

func (s *snapshotAge) record(snapshots []Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots = snapshots
}

func (s *snapshotAge) Collect(ch chan<- prometheus.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sn := range s.snapshots {
		ch <- prometheus.MustNewConstMetric(
			s.Desc,
			prometheus.CounterValue,
			sn.Age().Seconds(),
			sn.Name,
		)
	}
}

func (s *snapshotAge) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.Desc
}

type snapshotSize struct {
	mu        *sync.Mutex
	Desc      *prometheus.Desc
	snapshots []Snapshot
}

func newSnapshotSize() *snapshotSize {
	return &snapshotSize{
		mu: &sync.Mutex{},
		Desc: prometheus.NewDesc(
			"orderly_snapshots_size",
			"The size of the snapshots",
			[]string{"name"},
			nil,
		),
		snapshots: []Snapshot{},
	}
}

func (s *snapshotSize) record(snapshots []Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots = snapshots
}

func (s *snapshotSize) Collect(ch chan<- prometheus.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sn := range s.snapshots {
		ch <- prometheus.MustNewConstMetric(
			s.Desc,
			prometheus.CounterValue,
			float64(sn.Size),
			sn.Name,
		)
	}
}

func (s *snapshotSize) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.Desc
}
