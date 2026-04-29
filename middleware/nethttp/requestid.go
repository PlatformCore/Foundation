package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

const (
	// RequestIDKey is the context key for the request ID.
	RequestIDKey contextKey = "request_id"

	// HeaderXRequestID is the canonical request-ID header name.
	HeaderXRequestID = "X-Request-ID"

	// HeaderXCorrelationID is the correlation-ID header name (aliased).
	HeaderXCorrelationID = "X-Correlation-ID"
)

// RequestIDOptions configures the request ID middleware.
type RequestIDOptions struct {
	// Header is the header name to read/write the request ID. Defaults to X-Request-ID.
	Header string
	// Generator is a function that generates a new request ID. Defaults to UUID v4.
	Generator func() string
	// TrustIncoming controls whether an existing header value is trusted.
	TrustIncoming bool
}

// WithRequestID injects a unique request ID into every request context and response header.
// If the incoming request already carries the configured header and TrustIncoming is true,
// the existing value is reused.
func WithRequestID(opts ...RequestIDOptions) Middleware {
	cfg := RequestIDOptions{
		Header:        HeaderXRequestID,
		Generator:     func() string { return uuid.New().String() },
		TrustIncoming: true,
	}
	if len(opts) > 0 {
		o := opts[0]
		if o.Header != "" {
			cfg.Header = o.Header
		}
		if o.Generator != nil {
			cfg.Generator = o.Generator
		}
		cfg.TrustIncoming = o.TrustIncoming
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := ""
			if cfg.TrustIncoming {
				id = r.Header.Get(cfg.Header)
				if id == "" {
					id = r.Header.Get(HeaderXCorrelationID)
				}
			}
			if id == "" {
				id = cfg.Generator()
			}

			ctx := context.WithValue(r.Context(), RequestIDKey, id)
			w.Header().Set(cfg.Header, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestID retrieves the request ID from the context.
// Returns empty string if not set.
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return v
	}
	return ""
}
