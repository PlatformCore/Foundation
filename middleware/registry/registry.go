package registry

import (
	"fmt"
	"github.com/PlatformCore/libpackage/transport/core"
	"sync"
)

type Factory func(Config) (core.Middleware, error)
type Config map[string]any

type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

func New() *Registry { return &Registry{factories: map[string]Factory{}} }
func (r *Registry) Register(name string, f Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = f
}
func (r *Registry) Build(name string, cfg Config) (core.Middleware, error) {
	r.mu.RLock()
	f := r.factories[name]
	r.mu.RUnlock()
	if f == nil {
		return nil, fmt.Errorf("middleware factory %q not registered", name)
	}
	return f(cfg)
}
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.factories))
	for k := range r.factories {
		out = append(out, k)
	}
	return out
}

var Default = New()
