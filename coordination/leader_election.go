package coordination

import (
	"context"
	"errors"
	"sync"
	"time"
)

type LeadershipEvent struct {
	Leader  bool
	Lease   Lease
	At      time.Time
	Message string
}

type LeaderElector struct {
	store  LeaseStore
	key    string
	owner  string
	ttl    time.Duration
	renew  time.Duration
	events chan LeadershipEvent
	mu     sync.RWMutex
	lease  Lease
	leader bool
}

func NewLeaderElector(store LeaseStore, key, owner string, ttl time.Duration) *LeaderElector {
	if ttl <= 0 {
		ttl = 10 * time.Second
	}
	return &LeaderElector{
		store: store, key: key, owner: owner, ttl: ttl, renew: ttl / 3,
		events: make(chan LeadershipEvent, 16),
	}
}

func (e *LeaderElector) Events() <-chan LeadershipEvent { return e.events }

func (e *LeaderElector) IsLeader() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.leader
}

func (e *LeaderElector) Lease() Lease {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.lease
}

func (e *LeaderElector) Run(ctx context.Context) error {
	if e.store == nil {
		return errors.New("nil lease store")
	}
	ticker := time.NewTicker(e.renew)
	defer ticker.Stop()
	for {
		if err := e.step(ctx); err != nil {
			e.setLeader(false, Lease{}, err.Error())
		}
		select {
		case <-ctx.Done():
			e.resign(context.Background())
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (e *LeaderElector) step(ctx context.Context) error {
	e.mu.RLock()
	leader := e.leader
	lease := e.lease
	e.mu.RUnlock()
	if leader {
		next, err := e.store.Renew(ctx, lease, e.ttl)
		if err != nil {
			return err
		}
		e.setLeader(true, next, "renewed")
		return nil
	}
	l, err := e.store.TryAcquire(ctx, e.key, e.owner, e.ttl)
	if err != nil {
		return err
	}
	e.setLeader(true, l, "acquired")
	return nil
}

func (e *LeaderElector) resign(ctx context.Context) {
	e.mu.RLock()
	lease := e.lease
	leader := e.leader
	e.mu.RUnlock()
	if leader {
		_ = e.store.Release(ctx, lease)
		e.setLeader(false, Lease{}, "released")
	}
}

func (e *LeaderElector) setLeader(leader bool, lease Lease, msg string) {
	e.mu.Lock()
	changed := e.leader != leader || e.lease.Fencing != lease.Fencing
	e.leader = leader
	e.lease = lease
	e.mu.Unlock()
	if changed {
		select {
		case e.events <- LeadershipEvent{Leader: leader, Lease: lease, At: time.Now(), Message: msg}:
		default:
		}
	}
}
