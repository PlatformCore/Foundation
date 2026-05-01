package httpmw

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// CBState represents the circuit breaker state.
type CBState int

const (
	CBStateClosed   CBState = iota // Normal operation.
	CBStateOpen                    // Failing fast; no requests passed through.
	CBStateHalfOpen                // Probe requests allowed; evaluating recovery.
)

func (s CBState) String() string {
	switch s {
	case CBStateClosed:
		return "closed"
	case CBStateOpen:
		return "open"
	case CBStateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig configures the circuit breaker.
type CircuitBreakerConfig struct {
	// Name is used in logs and metrics. Required.
	Name string
	// MaxFailures is the number of consecutive failures before opening. Default: 5.
	MaxFailures int
	// OpenTimeout is how long the circuit stays open before entering half-open. Default: 10s.
	OpenTimeout time.Duration
	// HalfOpenMaxRequests is the maximum number of probe requests in half-open state. Default: 1.
	HalfOpenMaxRequests int
	// IsFailure determines whether a response counts as a failure.
	// If nil, any 5xx status code is a failure.
	IsFailure func(status int) bool
	// OnStateChange is called whenever the circuit breaker transitions state.
	OnStateChange func(name string, from, to CBState)
	// OnOpen is called when the circuit is open and a request is rejected.
	// If nil, a default 503 JSON response is written.
	OnOpen func(w http.ResponseWriter, r *http.Request)
}

// circuitBreaker is the state machine implementation.
type circuitBreaker struct {
	mu              sync.Mutex
	cfg             CircuitBreakerConfig
	state           CBState
	failures        int
	halfOpenPassed  int
	lastFailureTime time.Time
}

func newCircuitBreaker(cfg CircuitBreakerConfig) *circuitBreaker {
	if cfg.MaxFailures == 0 {
		cfg.MaxFailures = 5
	}
	if cfg.OpenTimeout == 0 {
		cfg.OpenTimeout = 10 * time.Second
	}
	if cfg.HalfOpenMaxRequests == 0 {
		cfg.HalfOpenMaxRequests = 1
	}
	if cfg.IsFailure == nil {
		cfg.IsFailure = func(status int) bool { return status >= 500 }
	}
	if cfg.OnOpen == nil {
		cfg.OnOpen = defaultCircuitOpen
	}
	return &circuitBreaker{cfg: cfg}
}

// allow checks whether the circuit breaker allows the request through.
// Returns false if the circuit is open.
func (cb *circuitBreaker) allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CBStateClosed:
		return true
	case CBStateOpen:
		if time.Since(cb.lastFailureTime) > cb.cfg.OpenTimeout {
			cb.transition(CBStateHalfOpen)
			cb.halfOpenPassed = 0
			return true
		}
		return false
	case CBStateHalfOpen:
		if cb.halfOpenPassed < cb.cfg.HalfOpenMaxRequests {
			cb.halfOpenPassed++
			return true
		}
		return false
	}
	return true
}

// record records the result of a request.
func (cb *circuitBreaker) record(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if success {
		switch cb.state {
		case CBStateHalfOpen:
			cb.transition(CBStateClosed)
			cb.failures = 0
		case CBStateClosed:
			cb.failures = 0
		}
		return
	}

	// Failure recorded.
	cb.lastFailureTime = time.Now()
	switch cb.state {
	case CBStateClosed:
		cb.failures++
		if cb.failures >= cb.cfg.MaxFailures {
			cb.transition(CBStateOpen)
		}
	case CBStateHalfOpen:
		cb.transition(CBStateOpen)
	}
}

func (cb *circuitBreaker) transition(to CBState) {
	from := cb.state
	cb.state = to
	if cb.cfg.OnStateChange != nil {
		go cb.cfg.OnStateChange(cb.cfg.Name, from, to)
	}
}

// WithCircuitBreaker returns a circuit breaker middleware.
// It monitors handler response status codes and opens the circuit on sustained failures.
func WithCircuitBreaker(cfg CircuitBreakerConfig) Middleware {
	cb := newCircuitBreaker(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cb.allow() {
				cfg.OnOpen(w, r)
				return
			}

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			success := !cfg.IsFailure(rw.status)
			cb.record(success)
		})
	}
}

func defaultCircuitOpen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Retry-After", "10")
	w.WriteHeader(http.StatusServiceUnavailable)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":      "service temporarily unavailable (circuit open)",
		"request_id": GetRequestID(r.Context()),
	})
}
