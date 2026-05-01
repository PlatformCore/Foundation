package httpmw

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// TimeoutConfig configures request timeout behavior.
type TimeoutConfig struct {
	// Timeout is the maximum duration allowed for a request. Required.
	Timeout time.Duration
	// SkipPaths lists paths exempt from timeout (e.g. long-polling endpoints).
	SkipPaths []string
	// OnTimeout is called when a request exceeds the timeout.
	OnTimeout func(w http.ResponseWriter, r *http.Request)
}

// WithTimeout returns a middleware that cancels the request context after the
// configured duration. Handlers must respect context cancellation.
func WithTimeout(cfg TimeoutConfig) Middleware {
	if cfg.OnTimeout == nil {
		cfg.OnTimeout = defaultTimeoutResponse
	}

	skipSet := make(map[string]struct{}, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skipSet[p] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, skip := skipSet[r.URL.Path]; skip {
				next.ServeHTTP(w, r)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout)
			defer cancel()

			tw := &timeoutWriter{
				ResponseWriter: w,
				done:           make(chan struct{}),
			}

			panicCh := make(chan interface{}, 1)
			go func() {
				defer func() {
					if rec := recover(); rec != nil {
						panicCh <- rec
					}
					close(tw.done)
				}()
				next.ServeHTTP(tw, r.WithContext(ctx))
			}()

			select {
			case <-ctx.Done():
				tw.mu.Lock()
				tw.timedOut = true
				tw.mu.Unlock()
				cfg.OnTimeout(w, r)
			case rec := <-panicCh:
				panic(rec)
			case <-tw.done:
				tw.mu.Lock()
				if tw.code != 0 {
					w.WriteHeader(tw.code)
				}
				if len(tw.buf) > 0 {
					_, _ = w.Write(tw.buf)
				}
				tw.mu.Unlock()
			}
		})
	}
}

type timeoutWriter struct {
	http.ResponseWriter
	mu       sync.Mutex
	buf      []byte
	code     int
	timedOut bool
	done     chan struct{}
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if !tw.timedOut {
		tw.code = code
	}
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, context.DeadlineExceeded
	}
	tw.buf = append(tw.buf, b...)
	return len(b), nil
}

func defaultTimeoutResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":      "request timeout",
		"request_id": GetRequestID(r.Context()),
	})
}
