package coordination

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrLockHeld     = errors.New("coordination: lock already held")
	ErrNotOwner     = errors.New("coordination: caller does not own lock")
	ErrLeaseExpired = errors.New("coordination: lease expired")
)

type Lease struct {
	Key       string
	Owner     string
	Fencing   uint64
	ExpiresAt time.Time
}

type LeaseStore interface {
	TryAcquire(ctx context.Context, key, owner string, ttl time.Duration) (Lease, error)
	Renew(ctx context.Context, lease Lease, ttl time.Duration) (Lease, error)
	Release(ctx context.Context, lease Lease) error
	Get(ctx context.Context, key string) (Lease, bool, error)
}

type MemoryLeaseStore struct {
	mu      sync.Mutex
	leases  map[string]Lease
	fencing map[string]uint64
	now     func() time.Time
}

func NewMemoryLeaseStore() *MemoryLeaseStore {
	return &MemoryLeaseStore{leases: make(map[string]Lease), fencing: make(map[string]uint64), now: time.Now}
}

func (s *MemoryLeaseStore) TryAcquire(ctx context.Context, key, owner string, ttl time.Duration) (Lease, error) {
	select {
	case <-ctx.Done():
		return Lease{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	if l, ok := s.leases[key]; ok && now.Before(l.ExpiresAt) {
		return Lease{}, ErrLockHeld
	}
	s.fencing[key]++
	l := Lease{Key: key, Owner: owner, Fencing: s.fencing[key], ExpiresAt: now.Add(ttl)}
	s.leases[key] = l
	return l, nil
}

func (s *MemoryLeaseStore) Renew(ctx context.Context, lease Lease, ttl time.Duration) (Lease, error) {
	select {
	case <-ctx.Done():
		return Lease{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.leases[lease.Key]
	if !ok || current.Owner != lease.Owner || current.Fencing != lease.Fencing {
		return Lease{}, ErrNotOwner
	}
	if s.now().After(current.ExpiresAt) {
		delete(s.leases, lease.Key)
		return Lease{}, ErrLeaseExpired
	}
	current.ExpiresAt = s.now().Add(ttl)
	s.leases[lease.Key] = current
	return current, nil
}

func (s *MemoryLeaseStore) Release(ctx context.Context, lease Lease) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.leases[lease.Key]
	if !ok {
		return nil
	}
	if current.Owner != lease.Owner || current.Fencing != lease.Fencing {
		return ErrNotOwner
	}
	delete(s.leases, lease.Key)
	return nil
}

func (s *MemoryLeaseStore) Get(ctx context.Context, key string) (Lease, bool, error) {
	select {
	case <-ctx.Done():
		return Lease{}, false, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.leases[key]
	if !ok {
		return Lease{}, false, nil
	}
	if s.now().After(l.ExpiresAt) {
		delete(s.leases, key)
		return Lease{}, false, nil
	}
	return l, true, nil
}
