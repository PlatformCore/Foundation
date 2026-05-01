package registry

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

type Kind string

const (
	KindHTTP       Kind = "http"
	KindGRPC       Kind = "grpc"
	KindEvent      Kind = "event"
	KindMiddleware Kind = "middleware"
	KindWorker     Kind = "worker"
	KindStorage    Kind = "storage"
)

type HealthState string

const (
	HealthUnknown   HealthState = "unknown"
	HealthHealthy   HealthState = "healthy"
	HealthDegraded  HealthState = "degraded"
	HealthUnhealthy HealthState = "unhealthy"
)

type Component struct {
	Name        string
	Kind        Kind
	Version     string
	Tags        map[string]string
	HealthCheck func(context.Context) error
	OnStart     func(context.Context) error
	OnStop      func(context.Context) error

	mu        sync.RWMutex
	state     HealthState
	lastError error
	lastCheck time.Time
}

func (c *Component) Health() HealthState { c.mu.RLock(); defer c.mu.RUnlock(); return c.state }
func (c *Component) LastError() error    { c.mu.RLock(); defer c.mu.RUnlock(); return c.lastError }
func (c *Component) Snapshot() ComponentSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return ComponentSnapshot{Name: c.Name, Kind: c.Kind, Version: c.Version, Tags: copyTags(c.Tags), State: c.state, LastError: c.lastError, LastCheck: c.lastCheck}
}

type ComponentSnapshot struct {
	Name      string
	Kind      Kind
	Version   string
	Tags      map[string]string
	State     HealthState
	LastError error
	LastCheck time.Time
}

type LifecycleHook func(event string, snap ComponentSnapshot)

type ComponentRegistry struct {
	mu       sync.RWMutex
	items    map[string]*Component
	hooks    []LifecycleHook
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

type ComponentOption func(*ComponentRegistry)

func WithComponentHealthInterval(d time.Duration) ComponentOption {
	return func(r *ComponentRegistry) {
		if d > 0 {
			r.interval = d
		}
	}
}
func WithComponentLifecycleHook(h LifecycleHook) ComponentOption {
	return func(r *ComponentRegistry) {
		if h != nil {
			r.hooks = append(r.hooks, h)
		}
	}
}

func NewComponentRegistry(opts ...ComponentOption) *ComponentRegistry {
	ctx, cancel := context.WithCancel(context.Background())
	r := &ComponentRegistry{items: map[string]*Component{}, interval: 30 * time.Second, ctx: ctx, cancel: cancel}
	for _, opt := range opts {
		opt(r)
	}
	go r.pollHealth()
	return r
}

func (r *ComponentRegistry) RegisterComponent(c *Component) error {
	if c == nil || c.Name == "" {
		return errors.New("component name is required")
	}
	c.mu.Lock()
	if c.state == "" {
		c.state = HealthUnknown
	}
	c.mu.Unlock()
	r.mu.Lock()
	r.items[c.Name] = c
	r.mu.Unlock()
	r.fireHooks("registered", c.Snapshot())
	return nil
}

func (r *ComponentRegistry) MustRegisterComponent(c *Component) {
	if err := r.RegisterComponent(c); err != nil {
		panic(err)
	}
}
func (r *ComponentRegistry) DeregisterComponent(name string) error {
	r.mu.Lock()
	c, ok := r.items[name]
	if ok {
		delete(r.items, name)
	}
	r.mu.Unlock()
	if !ok {
		return errors.New("component not found")
	}
	r.fireHooks("deregistered", c.Snapshot())
	return nil
}
func (r *ComponentRegistry) GetComponent(name string) (*Component, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.items[name]
	return c, ok
}
func (r *ComponentRegistry) GetComponentsByKind(kind Kind) []*Component {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []*Component{}
	for _, c := range r.items {
		if c.Kind == kind {
			out = append(out, c)
		}
	}
	return out
}
func (r *ComponentRegistry) GetComponentsByTag(key, value string) []*Component {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []*Component{}
	for _, c := range r.items {
		if c.Tags != nil && c.Tags[key] == value {
			out = append(out, c)
		}
	}
	return out
}
func (r *ComponentRegistry) ComponentSnapshots() []ComponentSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ComponentSnapshot, 0, len(r.items))
	for _, c := range r.items {
		out = append(out, c.Snapshot())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func (r *ComponentRegistry) IsHealthy() bool {
	for _, s := range r.ComponentSnapshots() {
		if s.State == HealthUnhealthy {
			return false
		}
	}
	return true
}
func (r *ComponentRegistry) ShutdownComponents(ctx context.Context) {
	r.cancel()
	for _, c := range r.ComponentSnapshots() {
		if comp, ok := r.GetComponent(c.Name); ok && comp.OnStop != nil {
			_ = comp.OnStop(ctx)
		}
	}
}

func (r *ComponentRegistry) pollHealth() {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			r.checkAll()
		}
	}
}
func (r *ComponentRegistry) checkAll() {
	r.mu.RLock()
	comps := make([]*Component, 0, len(r.items))
	for _, c := range r.items {
		comps = append(comps, c)
	}
	r.mu.RUnlock()
	for _, c := range comps {
		r.checkHealth(c)
	}
}
func (r *ComponentRegistry) checkHealth(c *Component) {
	if c.HealthCheck == nil {
		return
	}
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()
	err := c.HealthCheck(ctx)
	c.mu.Lock()
	c.lastCheck = time.Now()
	c.lastError = err
	if err != nil {
		c.state = HealthUnhealthy
	} else {
		c.state = HealthHealthy
	}
	snap := ComponentSnapshot{Name: c.Name, Kind: c.Kind, Version: c.Version, Tags: copyTags(c.Tags), State: c.state, LastError: c.lastError, LastCheck: c.lastCheck}
	c.mu.Unlock()
	r.fireHooks("health_checked", snap)
}
func (r *ComponentRegistry) fireHooks(event string, snap ComponentSnapshot) {
	for _, h := range r.hooks {
		h(event, snap)
	}
}
func copyTags(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
