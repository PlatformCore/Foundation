package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("cache: key not found")
	ErrClosed   = errors.New("cache: store closed")
)

type Entry[T any] struct {
	Key       string
	Value     T
	ExpiresAt time.Time
	CreatedAt time.Time
	Version   uint64
}

func (e Entry[T]) Expired(now time.Time) bool {
	return !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt)
}

type Store[T any] interface {
	Get(ctx context.Context, key string) (T, bool, error)
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Close() error
}

type Loader[T any] func(context.Context, string) (T, time.Duration, error)

type Metrics struct {
	Hits         uint64
	Misses       uint64
	Sets         uint64
	Deletes      uint64
	Evictions    uint64
	Refreshes    uint64
	LoadErrors   uint64
	InflightLoad int64
}

type Options struct {
	DefaultTTL        time.Duration
	CleanupInterval   time.Duration
	MaxEntries        int
	StaleWhileRefresh time.Duration
	NegativeTTL       time.Duration
}

func DefaultOptions() Options {
	return Options{
		DefaultTTL:        5 * time.Minute,
		CleanupInterval:   time.Minute,
		MaxEntries:        10000,
		StaleWhileRefresh: 30 * time.Second,
		NegativeTTL:       5 * time.Second,
	}
}

type SingleFlight[T any] struct {
	mu sync.Mutex
	m  map[string]*call[T]
}

type call[T any] struct {
	wg  sync.WaitGroup
	val T
	ttl time.Duration
	err error
}

func NewSingleFlight[T any]() *SingleFlight[T] {
	return &SingleFlight[T]{m: make(map[string]*call[T])}
}

func (g *SingleFlight[T]) Do(key string, fn func() (T, time.Duration, error)) (T, time.Duration, error, bool) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.ttl, c.err, true
	}
	c := new(call[T])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.ttl, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.ttl, c.err, false
}
