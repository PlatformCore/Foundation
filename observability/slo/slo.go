package slo

import (
	"math"
	"sync"
	"time"
)

type Window struct {
	Start    time.Time
	End      time.Time
	Good     uint64
	Total    uint64
	Budget   float64
	BurnRate float64
}

type Tracker struct {
	mu         sync.Mutex
	target     float64
	window     time.Duration
	buckets    []bucket
	bucketSize time.Duration
}

type bucket struct {
	start time.Time
	good  uint64
	total uint64
}

func NewTracker(target float64, window, bucketSize time.Duration) *Tracker {
	if target <= 0 || target >= 1 {
		target = 0.999
	}
	if window <= 0 {
		window = time.Hour
	}
	if bucketSize <= 0 {
		bucketSize = time.Minute
	}
	n := int(math.Ceil(float64(window) / float64(bucketSize)))
	return &Tracker{target: target, window: window, bucketSize: bucketSize, buckets: make([]bucket, n)}
}

func (t *Tracker) Observe(ok bool, at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	idx := int(at.UnixNano()/int64(t.bucketSize)) % len(t.buckets)
	start := at.Truncate(t.bucketSize)
	if !t.buckets[idx].start.Equal(start) {
		t.buckets[idx] = bucket{start: start}
	}
	t.buckets[idx].total++
	if ok {
		t.buckets[idx].good++
	}
}

func (t *Tracker) Snapshot(now time.Time) Window {
	if now.IsZero() {
		now = time.Now()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	start := now.Add(-t.window)
	var good, total uint64
	for _, b := range t.buckets {
		if b.start.After(start) || b.start.Equal(start) {
			good += b.good
			total += b.total
		}
	}
	availability := 1.0
	if total > 0 {
		availability = float64(good) / float64(total)
	}
	budget := (availability - t.target) / (1 - t.target)
	if budget < 0 {
		budget = 0
	}
	burn := 0.0
	if total > 0 {
		errRate := 1 - availability
		allowed := 1 - t.target
		if allowed > 0 {
			burn = errRate / allowed
		}
	}
	return Window{Start: start, End: now, Good: good, Total: total, Budget: budget, BurnRate: burn}
}
