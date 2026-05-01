package cache

import (
	"context"
	"sync/atomic"
	"time"
)

type ReadThrough[T any] struct {
	store   Store[T]
	loader  Loader[T]
	group   *SingleFlight[T]
	opts    Options
	metrics Metrics
}

func NewReadThrough[T any](store Store[T], loader Loader[T], opts Options) *ReadThrough[T] {
	if opts.DefaultTTL == 0 {
		opts = DefaultOptions()
	}
	return &ReadThrough[T]{store: store, loader: loader, group: NewSingleFlight[T](), opts: opts}
}

func (c *ReadThrough[T]) Get(ctx context.Context, key string) (T, error) {
	if v, ok, err := c.store.Get(ctx, key); err != nil {
		var zero T
		return zero, err
	} else if ok {
		atomic.AddUint64(&c.metrics.Hits, 1)
		return v, nil
	}
	atomic.AddUint64(&c.metrics.Misses, 1)
	atomic.AddInt64(&c.metrics.InflightLoad, 1)
	v, ttl, err, shared := c.group.Do(key, func() (T, time.Duration, error) {
		return c.loader(ctx, key)
	})
	atomic.AddInt64(&c.metrics.InflightLoad, -1)
	if err != nil {
		atomic.AddUint64(&c.metrics.LoadErrors, 1)
		var zero T
		return zero, err
	}
	if ttl == 0 {
		ttl = c.opts.DefaultTTL
	}
	if err := c.store.Set(ctx, key, v, ttl); err != nil {
		return v, err
	}
	if !shared {
		atomic.AddUint64(&c.metrics.Refreshes, 1)
	}
	return v, nil
}

func (c *ReadThrough[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	return c.store.Set(ctx, key, value, ttl)
}

func (c *ReadThrough[T]) Delete(ctx context.Context, key string) error {
	return c.store.Delete(ctx, key)
}

func (c *ReadThrough[T]) Metrics() Metrics {
	return c.metrics
}
