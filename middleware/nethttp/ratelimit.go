package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitStrategy defines how the rate limit key is extracted.
type RateLimitStrategy string

const (
	// StrategyIP limits by client IP address.
	StrategyIP RateLimitStrategy = "ip"
	// StrategyUser limits by authenticated user ID (requires WithAuth before this).
	StrategyUser RateLimitStrategy = "user"
	// StrategyTenant limits by tenant ID from JWT claims.
	StrategyTenant RateLimitStrategy = "tenant"
	// StrategyCustom uses a user-supplied key extractor.
	StrategyCustom RateLimitStrategy = "custom"
)

// RateLimitConfig configures the rate limiter.
type RateLimitConfig struct {
	// RequestsPerSecond is the sustained rate (tokens per second). Required.
	RequestsPerSecond float64
	// Burst is the maximum burst size. Defaults to RequestsPerSecond.
	Burst int
	// Strategy defines how limit keys are derived. Default: StrategyIP.
	Strategy RateLimitStrategy
	// KeyExtractor is used when Strategy is StrategyCustom.
	KeyExtractor func(r *http.Request) string
	// TTL is how long an idle limiter is kept alive. Default: 10m.
	TTL time.Duration
	// OnLimitExceeded is called when a request is rate-limited.
	// If nil, a default 429 JSON response is written.
	OnLimitExceeded func(w http.ResponseWriter, r *http.Request)
}

// rateLimiterEntry wraps a rate.Limiter with a last-seen timestamp.
type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// rateLimiterStore is a thread-safe store of per-key rate limiters.
type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*rateLimiterEntry
	rps      rate.Limit
	burst    int
	ttl      time.Duration
}

func newRateLimiterStore(rps float64, burst int, ttl time.Duration) *rateLimiterStore {
	s := &rateLimiterStore{
		limiters: make(map[string]*rateLimiterEntry),
		rps:      rate.Limit(rps),
		burst:    burst,
		ttl:      ttl,
	}
	go s.cleanup()
	return s
}

func (s *rateLimiterStore) get(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.limiters[key]
	if !ok {
		e = &rateLimiterEntry{limiter: rate.NewLimiter(s.rps, s.burst)}
		s.limiters[key] = e
	}
	e.lastSeen = time.Now()
	return e.limiter
}

func (s *rateLimiterStore) cleanup() {
	ticker := time.NewTicker(s.ttl / 2)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		cutoff := time.Now().Add(-s.ttl)
		for k, e := range s.limiters {
			if e.lastSeen.Before(cutoff) {
				delete(s.limiters, k)
			}
		}
		s.mu.Unlock()
	}
}

// WithRateLimit returns a token-bucket rate limiting middleware.
// Each unique key (IP, user, tenant, or custom) gets its own limiter.
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
	if cfg.OnLimitExceeded == nil {
		cfg.OnLimitExceeded = defaultTooManyRequests
	}

	store := newRateLimiterStore(cfg.RequestsPerSecond, cfg.Burst, cfg.TTL)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := extractRateLimitKey(r, cfg)
			limiter := store.get(key)

			if !limiter.Allow() {
				cfg.OnLimitExceeded(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func extractRateLimitKey(r *http.Request, cfg RateLimitConfig) string {
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
	w.Header().Set("Retry-After", "1")
	w.WriteHeader(http.StatusTooManyRequests)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":      "rate limit exceeded",
		"request_id": GetRequestID(r.Context()),
	})
}
