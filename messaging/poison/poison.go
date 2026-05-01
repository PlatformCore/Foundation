package poison

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type Decision string

const (
	DecisionProcess    Decision = "process"
	DecisionQuarantine Decision = "quarantine"
	DecisionDrop       Decision = "drop"
)

type Classifier struct {
	mu        sync.Mutex
	failures  map[string]failure
	threshold int
	window    time.Duration
}

type failure struct {
	count int
	first time.Time
	last  time.Time
}

func NewClassifier(threshold int, window time.Duration) *Classifier {
	if threshold <= 0 {
		threshold = 5
	}
	if window <= 0 {
		window = 10 * time.Minute
	}
	return &Classifier{failures: map[string]failure{}, threshold: threshold, window: window}
}

func Fingerprint(subject string, payload []byte) string {
	h := sha256.New()
	h.Write([]byte(subject))
	h.Write([]byte{0})
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Classifier) RecordFailure(ctx context.Context, fp string) Decision {
	select {
	case <-ctx.Done():
		return DecisionProcess
	default:
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	f := c.failures[fp]
	if f.first.IsZero() || now.Sub(f.first) > c.window {
		f = failure{first: now}
	}
	f.count++
	f.last = now
	c.failures[fp] = f
	if f.count >= c.threshold {
		return DecisionQuarantine
	}
	return DecisionProcess
}

func (c *Classifier) RecordSuccess(fp string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.failures, fp)
}

func (c *Classifier) Snapshot() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int, len(c.failures))
	for k, v := range c.failures {
		out[k] = v.count
	}
	return out
}
