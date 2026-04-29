package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig configures the distributed tracing middleware.
type TracingConfig struct {
	// ServiceName is the name of this service reported in spans. Required.
	ServiceName string
	// Tracer is a custom OpenTelemetry Tracer. If nil, the global tracer is used.
	Tracer trace.Tracer
	// Propagator is the TextMapPropagator used to extract/inject trace context.
	// If nil, the global propagator is used.
	Propagator propagation.TextMapPropagator
	// SpanNameFormatter formats the span name from the request.
	// If nil, "HTTP METHOD path" is used.
	SpanNameFormatter func(r *http.Request) string
	// SkipPaths lists paths that will not generate spans (e.g. health checks).
	SkipPaths []string
	// AdditionalAttributes allows injecting extra span attributes per-request.
	AdditionalAttributes func(r *http.Request) []attribute.KeyValue
}

// WithTracing returns an OpenTelemetry tracing middleware.
// It extracts existing trace context from incoming headers (W3C TraceContext / B3),
// starts a new server span, and propagates the context downstream.
func WithTracing(cfg TracingConfig) Middleware {
	tracer := cfg.Tracer
	if tracer == nil {
		tracer = otel.Tracer(cfg.ServiceName)
	}
	propagator := cfg.Propagator
	if propagator == nil {
		propagator = otel.GetTextMapPropagator()
	}
	spanNameFn := cfg.SpanNameFormatter
	if spanNameFn == nil {
		spanNameFn = func(r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}
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

			// Extract remote span context from headers.
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			spanName := spanNameFn(r)
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPMethod(r.Method),
					semconv.HTTPURL(r.URL.String()),
					semconv.HTTPScheme(scheme(r)),
					semconv.NetHostName(r.Host),
					semconv.HTTPUserAgent(r.UserAgent()),
					attribute.String("request_id", GetRequestID(ctx)),
				),
			)
			defer span.End()

			if cfg.AdditionalAttributes != nil {
				span.SetAttributes(cfg.AdditionalAttributes(r)...)
			}

			// Inject trace context into response headers for client visibility.
			propagator.Inject(ctx, propagation.HeaderCarrier(w.Header()))

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Record HTTP status on span.
			span.SetAttributes(semconv.HTTPStatusCode(rw.status))
			if rw.status >= 500 {
				span.SetStatus(codes.Error, http.StatusText(rw.status))
			} else {
				span.SetStatus(codes.Ok, "")
			}
		})
	}
}

func scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if s := r.Header.Get("X-Forwarded-Proto"); s != "" {
		return s
	}
	return "http"
}
