package healthcheck

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

type Status string

const (
	StatusUp       Status = "up"
	StatusDown     Status = "down"
	StatusDegraded Status = "degraded"
)

type DependencyCheck func(context.Context) error

type DependencyResult struct {
	Name      string
	Status    Status
	Latency   time.Duration
	Error     string
	CheckedAt time.Time
	Critical  bool
}

type DependencyRegistry struct {
	mu      sync.RWMutex
	checks  map[string]registeredCheck
	timeout time.Duration
}

type registeredCheck struct {
	fn       DependencyCheck
	critical bool
}

func NewDependencyRegistry(timeout time.Duration) *DependencyRegistry {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &DependencyRegistry{checks: make(map[string]registeredCheck), timeout: timeout}
}

func (r *DependencyRegistry) Register(name string, critical bool, check DependencyCheck) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks[name] = registeredCheck{fn: check, critical: critical}
}

func (r *DependencyRegistry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.checks, name)
}

func (r *DependencyRegistry) Check(ctx context.Context) []DependencyResult {
	r.mu.RLock()
	names := make([]string, 0, len(r.checks))
	snapshot := make(map[string]registeredCheck, len(r.checks))
	for name, check := range r.checks {
		names = append(names, name)
		snapshot[name] = check
	}
	r.mu.RUnlock()
	sort.Strings(names)

	results := make([]DependencyResult, len(names))
	var wg sync.WaitGroup
	for i, name := range names {
		i, name := i, name
		check := snapshot[name]
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			cctx, cancel := context.WithTimeout(ctx, r.timeout)
			defer cancel()
			err := check.fn(cctx)
			res := DependencyResult{Name: name, CheckedAt: time.Now(), Latency: time.Since(start), Critical: check.critical}
			if err != nil {
				res.Status = StatusDown
				res.Error = err.Error()
			} else {
				res.Status = StatusUp
			}
			results[i] = res
		}()
	}
	wg.Wait()
	return results
}

func Overall(results []DependencyResult) Status {
	degraded := false
	for _, r := range results {
		if r.Status == StatusDown && r.Critical {
			return StatusDown
		}
		if r.Status != StatusUp {
			degraded = true
		}
	}
	if degraded {
		return StatusDegraded
	}
	return StatusUp
}

func PingCheck(ping func(context.Context) error) DependencyCheck {
	return func(ctx context.Context) error {
		if ping == nil {
			return errors.New("nil ping function")
		}
		return ping(ctx)
	}
}
