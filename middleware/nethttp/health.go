package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// HealthStatus represents the health state of a component.
type HealthStatus string

const (
	HealthStatusUp   HealthStatus = "up"
	HealthStatusDown HealthStatus = "down"
)

// HealthCheck is a function that checks whether a dependency is healthy.
type HealthCheck func(ctx context.Context) error

// HealthConfig configures the health check endpoints.
type HealthConfig struct {
	// LivenessPath is the path for the liveness probe. Default: /healthz.
	LivenessPath string
	// ReadinessPath is the path for the readiness probe. Default: /readyz.
	ReadinessPath string
	// Checks maps component names to their check functions.
	Checks map[string]HealthCheck
	// CheckTimeout is the maximum time to wait for each check. Default: 5s.
	CheckTimeout time.Duration
}

// healthResponse is the JSON structure returned by health endpoints.
type healthResponse struct {
	Status     HealthStatus            `json:"status"`
	Components map[string]componentHealth `json:"components,omitempty"`
	Timestamp  time.Time               `json:"timestamp"`
}

type componentHealth struct {
	Status  HealthStatus `json:"status"`
	Message string       `json:"message,omitempty"`
}

// WithHealth registers /healthz (liveness) and /readyz (readiness) endpoints
// on the given ServeMux (or http.DefaultServeMux if nil).
// These paths are injected before the middleware chain so they always respond.
func WithHealth(cfg HealthConfig, mux *http.ServeMux) {
	if cfg.LivenessPath == "" {
		cfg.LivenessPath = "/healthz"
	}
	if cfg.ReadinessPath == "" {
		cfg.ReadinessPath = "/readyz"
	}
	if cfg.CheckTimeout == 0 {
		cfg.CheckTimeout = 5 * time.Second
	}
	if mux == nil {
		mux = http.DefaultServeMux
	}

	// Liveness: always returns 200 if the process is running.
	mux.HandleFunc(cfg.LivenessPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(healthResponse{
			Status:    HealthStatusUp,
			Timestamp: time.Now().UTC(),
		})
	})

	// Readiness: runs all registered health checks.
	mux.HandleFunc(cfg.ReadinessPath, func(w http.ResponseWriter, r *http.Request) {
		type result struct {
			name string
			err  error
		}

		ctx, cancel := context.WithTimeout(r.Context(), cfg.CheckTimeout)
		defer cancel()

		results := make(chan result, len(cfg.Checks))
		var wg sync.WaitGroup
		for name, check := range cfg.Checks {
			wg.Add(1)
			go func(n string, fn HealthCheck) {
				defer wg.Done()
				results <- result{name: n, err: fn(ctx)}
			}(name, check)
		}
		go func() {
			wg.Wait()
			close(results)
		}()

		resp := healthResponse{
			Status:     HealthStatusUp,
			Components: make(map[string]componentHealth, len(cfg.Checks)),
			Timestamp:  time.Now().UTC(),
		}
		for res := range results {
			if res.err != nil {
				resp.Status = HealthStatusDown
				resp.Components[res.name] = componentHealth{
					Status:  HealthStatusDown,
					Message: res.err.Error(),
				}
			} else {
				resp.Components[res.name] = componentHealth{Status: HealthStatusUp}
			}
		}

		statusCode := http.StatusOK
		if resp.Status == HealthStatusDown {
			statusCode = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(resp)
	})
}
