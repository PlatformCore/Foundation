package cache

import (
	"container/list"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type MemoryStore[T any] struct {
	mu      sync.RWMutex
	items   map[string]*list.Element
	lru     *list.List
	opts    Options
	closed  bool
	metrics Metrics
	stop    chan struct{}
}

type memoryNode[T any] struct {
	entry Entry[T]
}

func NewMemoryStore[T any](opts Options) *MemoryStore[T] {
	if opts.DefaultTTL == 0 {
		opts = DefaultOptions()
	}
	if opts.CleanupInterval <= 0 {
		opts.CleanupInterval = time.Minute
	}
	s := &MemoryStore[T]{
		items: make(map[string]*list.Element),
		lru:   list.New(),
		opts:  opts,
		stop:  make(chan struct{}),
	}
	go s.cleanupLoop()
	return s
}

func (s *MemoryStore[T]) Get(ctx context.Context, key string) (T, bool, error) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, false, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return zero, false, ErrClosed
	}
	el, ok := s.items[key]
	if !ok {
		atomic.AddUint64(&s.metrics.Misses, 1)
		return zero, false, nil
	}
	node := el.Value.(*memoryNode[T])
	if node.entry.Expired(time.Now()) {
		s.removeElement(el)
		atomic.AddUint64(&s.metrics.Misses, 1)
		return zero, false, nil
	}
	s.lru.MoveToFront(el)
	atomic.AddUint64(&s.metrics.Hits, 1)
	return node.entry.Value, true, nil
}

func (s *MemoryStore[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if ttl == 0 {
		ttl = s.opts.DefaultTTL
	}
	entry := Entry[T]{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if el, ok := s.items[key]; ok {
		node := el.Value.(*memoryNode[T])
		entry.Version = node.entry.Version + 1
		node.entry = entry
		s.lru.MoveToFront(el)
		atomic.AddUint64(&s.metrics.Sets, 1)
		return nil
	}
	el := s.lru.PushFront(&memoryNode[T]{entry: entry})
	s.items[key] = el
	atomic.AddUint64(&s.metrics.Sets, 1)
	s.enforceMaxLocked()
	return nil
}

func (s *MemoryStore[T]) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if el, ok := s.items[key]; ok {
		s.removeElement(el)
		atomic.AddUint64(&s.metrics.Deletes, 1)
	}
	return nil
}

func (s *MemoryStore[T]) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	s.items = make(map[string]*list.Element)
	s.lru.Init()
	return nil
}

func (s *MemoryStore[T]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	close(s.stop)
	s.items = nil
	s.lru.Init()
	return nil
}

func (s *MemoryStore[T]) Metrics() Metrics { return s.metrics }

func (s *MemoryStore[T]) cleanupLoop() {
	t := time.NewTicker(s.opts.CleanupInterval)
	defer t.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-t.C:
			s.PurgeExpired(time.Now())
		}
	}
}

func (s *MemoryStore[T]) PurgeExpired(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	removed := 0
	for key, el := range s.items {
		node := el.Value.(*memoryNode[T])
		if node.entry.Expired(now) {
			delete(s.items, key)
			s.lru.Remove(el)
			removed++
		}
	}
	if removed > 0 {
		atomic.AddUint64(&s.metrics.Evictions, uint64(removed))
	}
	return removed
}

func (s *MemoryStore[T]) enforceMaxLocked() {
	if s.opts.MaxEntries <= 0 {
		return
	}
	for s.lru.Len() > s.opts.MaxEntries {
		el := s.lru.Back()
		if el == nil {
			return
		}
		s.removeElement(el)
		atomic.AddUint64(&s.metrics.Evictions, 1)
	}
}

func (s *MemoryStore[T]) removeElement(el *list.Element) {
	node := el.Value.(*memoryNode[T])
	delete(s.items, node.entry.Key)
	s.lru.Remove(el)
}
