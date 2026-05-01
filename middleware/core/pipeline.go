// Package core is the canonical v11 middleware pipeline.
//
// It intentionally keeps middleware as a thin orchestration layer. Heavy logic
// belongs in resilience/* and observability/*, while transports only adapt their
// request/response objects into transport/core.Context.
package core

import (
	"sort"
	"strings"
	"sync"

	transportcore "github.com/PlatformCore/libpackage/transport/core"
)

// Stage defines the deterministic execution order for enterprise middleware.
type Stage int

const (
	StageIngress Stage = iota + 10
	StageIdentity
	StageSecurity
	StagePolicy
	StageResilience
	StageObservability
	StageBusiness
	StageEgress
)

// Spec is a named middleware registration. Name is used to prevent accidental
// duplicate responsibility, such as two retry layers or two timeout layers.
type Spec struct {
	Name     string
	Stage    Stage
	Priority int
	MW       transportcore.Middleware
	Replace  bool
}

// Registry stores canonical middleware by normalized name.
type Registry struct {
	mu    sync.RWMutex
	items map[string]Spec
	order []string
}

func NewRegistry() *Registry { return &Registry{items: map[string]Spec{}} }

func normalizeName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.ReplaceAll(name, "_", "-")
	return name
}

// Register adds a middleware once. If a middleware with the same name already
// exists, the new one is ignored unless Replace is true. This is the main guard
// against retry/timeout/ratelimit/observability duplication.
func (r *Registry) Register(spec Spec) *Registry {
	if r == nil || spec.MW == nil {
		return r
	}
	name := normalizeName(spec.Name)
	if name == "" {
		return r
	}
	if spec.Stage == 0 {
		spec.Stage = StageBusiness
	}
	spec.Name = name

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.items == nil {
		r.items = map[string]Spec{}
	}
	if _, exists := r.items[name]; exists && !spec.Replace {
		return r
	}
	if _, exists := r.items[name]; !exists {
		r.order = append(r.order, name)
	}
	r.items[name] = spec
	return r
}

func (r *Registry) Specs() []Spec {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Spec, 0, len(r.items))
	for _, name := range r.order {
		if spec, ok := r.items[name]; ok {
			out = append(out, spec)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Stage == out[j].Stage {
			if out[i].Priority == out[j].Priority {
				return out[i].Name < out[j].Name
			}
			return out[i].Priority < out[j].Priority
		}
		return out[i].Stage < out[j].Stage
	})
	return out
}

func (r *Registry) Middleware() []transportcore.Middleware {
	specs := r.Specs()
	out := make([]transportcore.Middleware, 0, len(specs))
	for _, spec := range specs {
		out = append(out, spec.MW)
	}
	return out
}

func (r *Registry) Then(h transportcore.Handler) transportcore.Handler {
	return transportcore.Chain(h, r.Middleware()...)
}

func Compose(specs ...Spec) transportcore.Middleware {
	reg := NewRegistry()
	for _, spec := range specs {
		reg.Register(spec)
	}
	return transportcore.Compose(reg.Middleware()...)
}

func Named(name string, stage Stage, mw transportcore.Middleware) Spec {
	return Spec{Name: name, Stage: stage, MW: mw}
}
