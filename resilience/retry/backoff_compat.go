package retry

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

const BackoffDecorrelated BackoffStrategy = 100

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Strategy    BackoffStrategy
	Jitter      float64
	Multiplier  float64
	RetryIf     func(attempt int, err error) bool
	OnRetry     func(attempt int, delay time.Duration, err error)
	OnSuccess   func(attempt int, totalDuration time.Duration)
	OnFailure   func(attempts int, totalDuration time.Duration, lastErr error)
}

type RetryStats struct {
	Attempts      int
	TotalDelay    time.Duration
	TotalDuration time.Duration
	LastErr       error
	Succeeded     bool
}

var ErrRetryNotAllowed = errors.New("retry not allowed for this error")

type Retryer struct{ cfg RetryConfig }

func NewRetryer(cfg RetryConfig) *Retryer { return &Retryer{cfg: cfg} }
func (r *Retryer) Do(ctx context.Context, fn func(ctx context.Context, attempt int) error) (*RetryStats, error) {
	start := time.Now()
	stats := &RetryStats{}
	cfg := DefaultConfig()
	if r.cfg.MaxAttempts != 0 {
		cfg.MaxAttempts = r.cfg.MaxAttempts
	}
	if r.cfg.BaseDelay > 0 {
		cfg.InitialDelay = r.cfg.BaseDelay
	}
	if r.cfg.MaxDelay > 0 {
		cfg.MaxDelay = r.cfg.MaxDelay
	}
	if r.cfg.Multiplier > 0 {
		cfg.Multiplier = r.cfg.Multiplier
	}
	if r.cfg.Jitter >= 0 {
		cfg.Jitter = r.cfg.Jitter
	}
	if r.cfg.Strategy != BackoffDecorrelated {
		cfg.Strategy = r.cfg.Strategy
	}
	if r.cfg.RetryIf != nil {
		cfg.RetryIf = func(err error) bool { return r.cfg.RetryIf(stats.Attempts, err) }
	}
	if r.cfg.OnRetry != nil {
		cfg.OnRetry = r.cfg.OnRetry
	}
	res := Do(ctx, cfg, func(c context.Context) error {
		stats.Attempts++
		return fn(c, stats.Attempts-1)
	})
	stats.TotalDuration = time.Since(start)
	stats.LastErr = res.Err
	stats.Succeeded = res.Err == nil
	if res.Err == nil && r.cfg.OnSuccess != nil {
		r.cfg.OnSuccess(stats.Attempts-1, stats.TotalDuration)
	}
	if res.Err != nil && r.cfg.OnFailure != nil {
		r.cfg.OnFailure(stats.Attempts, stats.TotalDuration, res.Err)
	}
	if res.Err != nil && cfg.RetryIf != nil && !cfg.RetryIf(res.Err) {
		return stats, fmt.Errorf("%w: %w", ErrRetryNotAllowed, res.Err)
	}
	return stats, res.Err
}
func ExponentialRetryer(maxAttempts int, base, maxDelay time.Duration) *Retryer {
	return NewRetryer(RetryConfig{MaxAttempts: maxAttempts, BaseDelay: base, MaxDelay: maxDelay, Strategy: BackoffExponential, Jitter: 0.25, Multiplier: 2.0})
}
func LinearRetryer(maxAttempts int, base, maxDelay time.Duration) *Retryer {
	return NewRetryer(RetryConfig{MaxAttempts: maxAttempts, BaseDelay: base, MaxDelay: maxDelay, Strategy: BackoffLinear, Jitter: 0.1, Multiplier: 1.0})
}

type RetryCounter struct{ totalCalls, totalRetries, totalFailures, totalSuccesses int64 }

func (c *RetryCounter) Record(stats *RetryStats) {
	if stats == nil {
		return
	}
	atomic.AddInt64(&c.totalCalls, 1)
	atomic.AddInt64(&c.totalRetries, int64(stats.Attempts-1))
	if stats.Succeeded {
		atomic.AddInt64(&c.totalSuccesses, 1)
	} else {
		atomic.AddInt64(&c.totalFailures, 1)
	}
}
func (c *RetryCounter) Snapshot() (calls, retries, failures, successes int64) {
	return atomic.LoadInt64(&c.totalCalls), atomic.LoadInt64(&c.totalRetries), atomic.LoadInt64(&c.totalFailures), atomic.LoadInt64(&c.totalSuccesses)
}
