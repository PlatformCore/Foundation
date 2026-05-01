package adaptive_limiter

import (
	"context"
	"math"
	"sync"
	"time"
)

type Sample struct {
	Latency  time.Duration
	Success  bool
	InFlight int
}

type Limiter struct {
	mu          sync.Mutex
	limit       int
	min         int
	max         int
	inflight    int
	ewmaLatency float64
	errorRate   float64
	alpha       float64
	lastAdjust  time.Time
	interval    time.Duration
}

type Option func(*Limiter)

func WithBounds(min, max int) Option { return func(l *Limiter) { l.min, l.max = min, max } }
func WithInterval(d time.Duration) Option {
	return func(l *Limiter) {
		if d > 0 {
			l.interval = d
		}
	}
}

func New(initial int, opts ...Option) *Limiter {
	if initial <= 0 {
		initial = 32
	}
	l := &Limiter{limit: initial, min: 1, max: 10000, alpha: 0.2, interval: time.Second, lastAdjust: time.Now()}
	for _, opt := range opts {
		opt(l)
	}
	if l.limit < l.min {
		l.limit = l.min
	}
	if l.limit > l.max {
		l.limit = l.max
	}
	return l
}

func (l *Limiter) Acquire(ctx context.Context) (func(Sample), error) {
	t := time.NewTicker(time.Millisecond)
	defer t.Stop()
	for {
		l.mu.Lock()
		if l.inflight < l.limit {
			l.inflight++
			startInFlight := l.inflight
			l.mu.Unlock()
			start := time.Now()
			return func(s Sample) {
				if s.Latency == 0 {
					s.Latency = time.Since(start)
				}
				if s.InFlight == 0 {
					s.InFlight = startInFlight
				}
				l.Release(s)
			}, nil
		}
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:
		}
	}
}

func (l *Limiter) Release(s Sample) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.inflight > 0 {
		l.inflight--
	}
	latMs := float64(s.Latency.Milliseconds())
	if l.ewmaLatency == 0 {
		l.ewmaLatency = latMs
	} else {
		l.ewmaLatency = l.alpha*latMs + (1-l.alpha)*l.ewmaLatency
	}
	errSample := 0.0
	if !s.Success {
		errSample = 1
	}
	l.errorRate = l.alpha*errSample + (1-l.alpha)*l.errorRate
	if time.Since(l.lastAdjust) < l.interval {
		return
	}
	l.lastAdjust = time.Now()

	overloaded := l.errorRate > 0.05 || (l.ewmaLatency > 0 && latMs > l.ewmaLatency*2)
	if overloaded {
		l.limit = int(math.Max(float64(l.min), float64(l.limit)*0.85))
	} else if s.InFlight >= int(float64(l.limit)*0.8) {
		l.limit = int(math.Min(float64(l.max), float64(l.limit+1)))
	}
}

func (l *Limiter) Limit() int    { l.mu.Lock(); defer l.mu.Unlock(); return l.limit }
func (l *Limiter) InFlight() int { l.mu.Lock(); defer l.mu.Unlock(); return l.inflight }
