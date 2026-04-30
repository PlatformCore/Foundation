// Package http provides enterprise HTTP middleware that integrates obslib's
// ID propagation, tracing, metrics, and structured logging in a single,
// composable middleware chain.
package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"

	"github.com/PlatformCore/libpackage/middleware/ids"
	"github.com/PlatformCore/libpackage/middleware/propagation"
	obttrace "github.com/PlatformCore/libpackage/middleware/trace"
)

// ─── Middleware Type ──────────────────────────────────────────────────────────

// Middleware is a standard net/http middleware function.
type Middleware func(http.Handler) http.Handler

// Chain composes middleware in order (first = outermost).
func Chain(mws ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}

// ─── responseWriter ───────────────────────────────────────────────────────────

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
	wrote  bool
}

func wrap(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wrote {
		rw.status = code
		rw.wrote = true
		rw.ResponseWriter.WriteHeader(code)
	}
}
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wrote {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}
func (rw *responseWriter) Unwrap() http.ResponseWriter { return rw.ResponseWriter }

// ─── Propagation Middleware ───────────────────────────────────────────────────

// WithPropagation is the foundational middleware. It MUST be the outermost
// middleware in any chain. It:
//   1. Extracts all IDs from incoming headers (W3C, B3, X-Obs-Context).
//   2. Generates missing IDs (RequestID, TraceID, SpanID).
//   3. Injects all IDs into the request context.
//   4. Stamps all IDs onto the response headers.
//   5. Generates a ResponseID and stamps it too.
func WithPropagation() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := propagation.HTTPExtract(r.Context(), r.Header)

			// Generate response ID.
			rid := ids.GetRequestID(ctx)
			responseID := ids.NewResponseID(rid)
			ctx = ids.WithResponseID(ctx, responseID)

			// Inject all IDs into response before handler runs.
			propagation.HTTPInject(ctx, w.Header())
			w.Header().Set(ids.HeaderResponseID, string(responseID))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── Observability Middleware ─────────────────────────────────────────────────

// ObsConfig configures the combined observability middleware.
type ObsConfig struct {
	// Tracer is the enterprise tracer. Required.
	Tracer *obttrace.Tracer
	// Logger for structured request logging. Required.
	Logger *zap.Logger
	// Metrics config for Prometheus. Optional.
	Metrics *MetricsConfig
	// SkipPaths are excluded from logging, tracing, and metrics.
	SkipPaths []string
	// SlowThreshold marks requests as slow. Default: 1s.
	SlowThreshold time.Duration
	// LogRequestHeaders logs selected headers. Empty = none.
	LogRequestHeaders []string
}

// MetricsConfig holds Prometheus metric registrations.
type MetricsConfig struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
	ResponseBytes    *prometheus.HistogramVec
}

// DefaultMetrics creates and registers standard Prometheus metrics.
func DefaultMetrics(reg prometheus.Registerer, namespace, subsystem string) *MetricsConfig {
	if namespace == "" {
		namespace = "obs"
	}
	if subsystem == "" {
		subsystem = "http"
	}
	m := &MetricsConfig{
		RequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Subsystem: subsystem, Name: "requests_total",
			Help: "Total HTTP requests.",
		}, []string{"method", "path", "status", "service"}),
		RequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace, Subsystem: subsystem, Name: "request_duration_seconds",
			Help:    "HTTP request latency.",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}, []string{"method", "path", "status", "service"}),
		RequestsInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace, Subsystem: subsystem, Name: "requests_in_flight",
			Help: "Active HTTP requests.",
		}),
		ResponseBytes: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace, Subsystem: subsystem, Name: "response_size_bytes",
			Help:    "HTTP response sizes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		}, []string{"method", "path", "service"}),
	}
	if reg != nil {
		reg.MustRegister(m.RequestsTotal, m.RequestDuration, m.RequestsInFlight, m.ResponseBytes)
	}
	return m
}

// WithObservability is the combined tracing + logging + metrics middleware.
// It emits a structured log line, a trace span, and Prometheus metrics for
// every request, enriched with all obslib IDs from the context.
func WithObservability(cfg ObsConfig) Middleware {
	if cfg.SlowThreshold == 0 {
		cfg.SlowThreshold = time.Second
	}
	skipSet := makeSkipSet(cfg.SkipPaths)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, skip := skipSet[r.URL.Path]; skip {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			// ── Tracing ─────────────────────────────────────────────────────
			spanName := r.Method + " " + r.URL.Path
			ctx, span := cfg.Tracer.Start(r.Context(), spanName)
			span.SetAttr(
				semconv.HTTPMethod(r.Method),
				semconv.HTTPURL(r.URL.String()),
				semconv.NetHostName(r.Host),
				semconv.HTTPUserAgent(r.UserAgent()),
			)

			// ── Metrics: in-flight ───────────────────────────────────────────
			if cfg.Metrics != nil {
				cfg.Metrics.RequestsInFlight.Inc()
				defer cfg.Metrics.RequestsInFlight.Dec()
			}

			rw := wrap(w)
			next.ServeHTTP(rw, r.WithContext(ctx))

			latency := time.Since(start)
			statusStr := strconv.Itoa(rw.status)
			carrier := ids.CarrierFromContext(ctx)

			// ── Span finalise ────────────────────────────────────────────────
			span.SetHTTP(r.Method, r.URL.Path, rw.status)
			if rw.status >= 500 {
				span.OTel().SetStatus(codes.Error, http.StatusText(rw.status))
			}
			span.OTel().SetAttributes(
				attribute.String("request.id", string(carrier.RequestID)),
				attribute.String("response.id", string(carrier.ResponseID)),
				attribute.String("correlation.id", string(carrier.CorrelationID)),
			)
			span.End()

			// ── Metrics ──────────────────────────────────────────────────────
			svcName := cfg.Tracer.Config().ServiceName
			if cfg.Metrics != nil {
				lbls := []string{r.Method, r.URL.Path, statusStr, svcName}
				cfg.Metrics.RequestsTotal.WithLabelValues(lbls...).Inc()
				cfg.Metrics.RequestDuration.WithLabelValues(lbls...).Observe(latency.Seconds())
				cfg.Metrics.ResponseBytes.WithLabelValues(r.Method, r.URL.Path, svcName).
					Observe(float64(rw.size))
			}

			// ── Logging ──────────────────────────────────────────────────────
			fields := []zap.Field{
				zap.String("request_id", string(carrier.RequestID)),
				zap.String("response_id", string(carrier.ResponseID)),
				zap.String("correlation_id", string(carrier.CorrelationID)),
				zap.String("trace_id", obttrace.TraceIDFromContext(ctx)),
				zap.String("span_id", obttrace.SpanIDFromContext(ctx)),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status", rw.status),
				zap.Int("bytes", rw.size),
				zap.Duration("latency", latency),
				zap.String("ip", realIP(r)),
				zap.String("user_agent", r.UserAgent()),
				zap.String("service", svcName),
			}
			for _, hdr := range cfg.LogRequestHeaders {
				fields = append(fields, zap.String("req."+hdr, r.Header.Get(hdr)))
			}
			if latency > cfg.SlowThreshold {
				fields = append(fields, zap.Bool("slow", true))
			}

			switch {
			case rw.status >= 500:
				cfg.Logger.Error("http request", fields...)
			case rw.status >= 400:
				cfg.Logger.Warn("http request", fields...)
			default:
				cfg.Logger.Info("http request", fields...)
			}
		})
	}
}

// ─── Recovery ─────────────────────────────────────────────────────────────────

// WithRecovery recovers panics, logs them with full obslib context, and
// returns a 500 JSON response.
func WithRecovery(log *zap.Logger, tracer *obttrace.Tracer) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					carrier := ids.CarrierFromContext(r.Context())
					log.Error("panic recovered",
						zap.Any("panic", rec),
						zap.String("request_id", string(carrier.RequestID)),
						zap.String("trace_id", obttrace.TraceIDFromContext(r.Context())),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)
					if span := obttrace.SpanFromContext(r.Context()); span != nil {
						span.RecordError(r.Context().Err())
					}
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]string{
						"error":      "internal server error",
						"request_id": string(carrier.RequestID),
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func makeSkipSet(paths []string) map[string]struct{} {
	s := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		s[p] = struct{}{}
	}
	return s
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		if i := len(ip); i > 0 {
			for j := 0; j < i; j++ {
				if ip[j] == ',' {
					return ip[:j]
				}
			}
		}
		return ip
	}
	return r.RemoteAddr
}
