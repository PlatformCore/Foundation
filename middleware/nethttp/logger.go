package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerOptions configures the logger middleware.
type LoggerOptions struct {
	// Logger is the zap logger instance. Required.
	Logger *zap.Logger
	// SlowRequestThreshold defines when a request is considered slow. Defaults to 1s.
	SlowRequestThreshold time.Duration
	// SkipPaths is a list of paths that will not be logged.
	SkipPaths []string
	// LogRequestBody enables request body logging (use carefully — impacts performance).
	LogRequestBody bool
	// FieldExtractors allows injecting extra zap fields per-request.
	FieldExtractors []func(r *http.Request) zap.Field
}

// responseWriter wraps http.ResponseWriter to capture status code and bytes written.
type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Unwrap supports http.ResponseController (Go 1.20+).
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// WithLogger returns a middleware that logs every HTTP request using structured zap logging.
// It captures method, path, status, latency, bytes written, request ID, and user-agent.
func WithLogger(logger *zap.Logger, opts ...LoggerOptions) Middleware {
	cfg := LoggerOptions{
		Logger:               logger,
		SlowRequestThreshold: time.Second,
	}
	if len(opts) > 0 {
		o := opts[0]
		if o.SlowRequestThreshold > 0 {
			cfg.SlowRequestThreshold = o.SlowRequestThreshold
		}
		if len(o.SkipPaths) > 0 {
			cfg.SkipPaths = o.SkipPaths
		}
		cfg.LogRequestBody = o.LogRequestBody
		cfg.FieldExtractors = o.FieldExtractors
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

			start := time.Now()
			rw := newResponseWriter(w)

			defer func() {
				latency := time.Since(start)
				requestID := GetRequestID(r.Context())

				fields := []zap.Field{
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("query", r.URL.RawQuery),
					zap.Int("status", rw.status),
					zap.Int("bytes", rw.size),
					zap.Duration("latency", latency),
					zap.String("ip", realIP(r)),
					zap.String("user_agent", r.UserAgent()),
					zap.String("request_id", requestID),
					zap.String("proto", r.Proto),
				}

				for _, extract := range cfg.FieldExtractors {
					fields = append(fields, extract(r))
				}

				level := zapcore.InfoLevel
				if rw.status >= 500 {
					level = zapcore.ErrorLevel
				} else if rw.status >= 400 {
					level = zapcore.WarnLevel
				}
				if latency > cfg.SlowRequestThreshold {
					fields = append(fields, zap.Bool("slow_request", true))
					level = zapcore.WarnLevel
				}

				if ce := cfg.Logger.Check(level, "http request"); ce != nil {
					ce.Write(fields...)
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

// realIP extracts the real client IP respecting common proxy headers.
func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For may be a comma-separated list; take the first.
		for i := 0; i < len(ip); i++ {
			if ip[i] == ',' {
				return ip[:i]
			}
		}
		return ip
	}
	return r.RemoteAddr
}
