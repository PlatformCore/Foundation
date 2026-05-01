package httpmw

import (
	"encoding/json"
	"net/http"
	"time"

	canonical "github.com/PlatformCore/libpackage/resilience/ratelimit"
)

// RateLimitStrategy defines how the rate limit key is extracted.
type RateLimitStrategy string

const (
	StrategyIP     RateLimitStrategy = "ip"
	StrategyUser   RateLimitStrategy = "user"
	StrategyTenant RateLimitStrategy = "tenant"
	StrategyCustom RateLimitStrategy = "custom"
)

// RateLimitConfig configures HTTP rate limiting.
// The middleware delegates all limiter logic to resilience/ratelimit.
type RateLimitConfig struct {
	RequestsPerSecond float64
	Burst             int
	Strategy          RateLimitStrategy
	KeyExtractor      func(r *http.Request) string
	TTL               time.Duration
	Namespace         string
	RouteFunc         func(r *http.Request) string
	Limiter           *canonical.Limiter
	Policy            canonical.Policy
	OnLimitExceeded   func(w http.ResponseWriter, r *http.Request)
}

// WithRateLimit returns an HTTP middleware backed by the canonical resilience/ratelimit package.
func WithRateLimit(cfg RateLimitConfig) Middleware {
	if cfg.Strategy == "" {
		cfg.Strategy = StrategyIP
	}
	if cfg.Burst == 0 {
		cfg.Burst = int(cfg.RequestsPerSecond)
		if cfg.Burst < 1 {
			cfg.Burst = 1
		}
	}
	if cfg.TTL == 0 {
		cfg.TTL = 10 * time.Minute
	}
	if cfg.Namespace == "" {
		cfg.Namespace = "http"
	}
	if cfg.OnLimitExceeded == nil {
		cfg.OnLimitExceeded = defaultTooManyRequests
	}

	limiter := cfg.Limiter
	if limiter == nil {
		limit := int64(cfg.Burst)
		if cfg.RequestsPerSecond > 0 {
			limit = int64(cfg.RequestsPerSecond)
		}
		if limit < 1 {
			limit = 1
		}
		policy := cfg.Policy.Normalize()
		policy.Name = firstNonEmpty(policy.Name, "http-default")
		if cfg.Policy.Limit == 0 {
			policy.Limit = limit
		}
		if cfg.Policy.Burst == 0 {
			policy.Burst = int64(cfg.Burst)
		}
		if cfg.Policy.Window == 0 {
			policy.Window = time.Second
		}
		if cfg.Policy.Strategy == "" {
			policy.Strategy = canonical.StrategyTokenBucket
		}
		limiter = canonical.New(canonical.Options{Policy: policy, Store: canonical.NewMemoryStore()})
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity := extractRateLimitIdentity(r, cfg)
			key := canonical.NewKey(cfg.Namespace, identity)
			key.Method = r.Method
			if cfg.RouteFunc != nil {
				key.Route = cfg.RouteFunc(r)
			} else {
				key.Route = r.URL.Path
			}
			if c := GetClaims(r.Context()); c != nil {
				key.Tenant = c.TenantID
			}

			result, err := limiter.Allow(r.Context(), key)
			if err != nil || !result.Allowed {
				if result.RetryAfter > 0 {
					w.Header().Set("Retry-After", result.RetryAfter.String())
				}
				canonical.WriteHeaders(w.Header(), result)
				cfg.OnLimitExceeded(w, r)
				return
			}
			canonical.WriteHeaders(w.Header(), result)
			next.ServeHTTP(w, r)
		})
	}
}

func extractRateLimitIdentity(r *http.Request, cfg RateLimitConfig) string {
	switch cfg.Strategy {
	case StrategyUser:
		if c := GetClaims(r.Context()); c != nil && c.UserID != "" {
			return "user:" + c.UserID
		}
		return "ip:" + realIP(r)
	case StrategyTenant:
		if c := GetClaims(r.Context()); c != nil && c.TenantID != "" {
			return "tenant:" + c.TenantID
		}
		return "ip:" + realIP(r)
	case StrategyCustom:
		if cfg.KeyExtractor != nil {
			return cfg.KeyExtractor(r)
		}
		fallthrough
	default:
		return "ip:" + realIP(r)
	}
}

func defaultTooManyRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if w.Header().Get("Retry-After") == "" {
		w.Header().Set("Retry-After", "1")
	}
	w.WriteHeader(http.StatusTooManyRequests)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded", "request_id": GetRequestID(r.Context())})
}

func firstNonEmpty(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}
