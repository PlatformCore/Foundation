// Package grpc provides enterprise gRPC unary and streaming interceptors
// with full obslib ID propagation, distributed tracing, Prometheus metrics,
// and structured logging.
package grpc

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/PlatformCore/Foundation/middleware/ids"
	"github.com/PlatformCore/Foundation/middleware/propagation"
	obttrace "github.com/PlatformCore/Foundation/middleware/trace"
)

// ─── Metrics ──────────────────────────────────────────────────────────────────

// GRPCMetrics holds Prometheus metrics for gRPC servers.
type GRPCMetrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	InFlight        prometheus.Gauge
}

// DefaultGRPCMetrics registers standard gRPC Prometheus metrics.
func DefaultGRPCMetrics(reg prometheus.Registerer, namespace string) *GRPCMetrics {
	if namespace == "" {
		namespace = "obs"
	}
	m := &GRPCMetrics{
		RequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Subsystem: "grpc", Name: "requests_total",
			Help: "Total gRPC requests.",
		}, []string{"grpc_service", "grpc_method", "grpc_code", "service"}),
		RequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace, Subsystem: "grpc", Name: "request_duration_seconds",
			Help:    "gRPC request latency.",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		}, []string{"grpc_service", "grpc_method", "grpc_code", "service"}),
		InFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace, Subsystem: "grpc", Name: "requests_in_flight",
			Help: "Active gRPC requests.",
		}),
	}
	if reg != nil {
		reg.MustRegister(m.RequestsTotal, m.RequestDuration, m.InFlight)
	}
	return m
}

// ─── InterceptorConfig ────────────────────────────────────────────────────────

// InterceptorConfig configures gRPC server interceptors.
type InterceptorConfig struct {
	Tracer      *obttrace.Tracer
	Logger      *zap.Logger
	Metrics     *GRPCMetrics
	SkipMethods []string // full gRPC method paths to skip, e.g. "/grpc.health.v1.Health/Check"
}

// ─── Unary Server Interceptor ────────────────────────────────────────────────

// UnaryServerInterceptor returns a gRPC unary server interceptor that:
//   1. Extracts obslib IDs from incoming gRPC metadata.
//   2. Starts an enterprise span for the call.
//   3. Emits Prometheus metrics.
//   4. Writes a structured log line.
//   5. Recovers from panics.
//   6. Injects response IDs into trailing metadata.
func UnaryServerInterceptor(cfg InterceptorConfig) grpc.UnaryServerInterceptor {
	skipSet := makeSkipSet(cfg.SkipMethods)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, skip := skipSet[info.FullMethod]; skip {
			return handler(ctx, req)
		}

		// ── Propagation ──────────────────────────────────────────────────────
		ctx = propagation.GRPCExtract(ctx)
		carrier := ids.CarrierFromContext(ctx)
		responseID := ids.NewResponseID(carrier.RequestID)

		// ── Tracing ──────────────────────────────────────────────────────────
		grpcSvc, grpcMethod := splitMethod(info.FullMethod)
		spanName := info.FullMethod
		ctx, span := cfg.Tracer.Start(ctx, spanName)
		span.SetGRPC(grpcSvc, grpcMethod, 0)
		span.SetAttr(
			attribute.String("request.id", string(carrier.RequestID)),
			attribute.String("response.id", string(responseID)),
		)

		// ── Metrics: in-flight ───────────────────────────────────────────────
		if cfg.Metrics != nil {
			cfg.Metrics.InFlight.Inc()
			defer cfg.Metrics.InFlight.Dec()
		}

		start := time.Now()

		// ── Panic recovery ───────────────────────────────────────────────────
		defer func() {
			if rec := recover(); rec != nil {
				cfg.Logger.Error("grpc panic recovered",
					zap.Any("panic", rec),
					zap.String("method", info.FullMethod),
					zap.String("request_id", string(carrier.RequestID)),
					zap.String("trace_id", obttrace.TraceIDFromContext(ctx)),
				)
				span.RecordError(status.Errorf(codes.Internal, "panic: %v", rec))
				span.End()
			}
		}()

		// ── Call handler ────────────────────────────────────────────────────
		resp, err := handler(ctx, req)
		latency := time.Since(start)
		grpcCode := statusCode(err)
		codeStr := grpcCode.String()

		// ── Finalise span ────────────────────────────────────────────────────
		span.SetGRPC(grpcSvc, grpcMethod, uint32(grpcCode))
		if err != nil {
			span.RecordError(err)
		}
		span.End()

		// ── Metrics ──────────────────────────────────────────────────────────
		svcName := cfg.Tracer.Config().ServiceName
		if cfg.Metrics != nil {
			lbls := []string{grpcSvc, grpcMethod, codeStr, svcName}
			cfg.Metrics.RequestsTotal.WithLabelValues(lbls...).Inc()
			cfg.Metrics.RequestDuration.WithLabelValues(lbls...).Observe(latency.Seconds())
		}

		// ── Logging ──────────────────────────────────────────────────────────
		fields := []zap.Field{
			zap.String("request_id", string(carrier.RequestID)),
			zap.String("response_id", string(responseID)),
			zap.String("correlation_id", string(carrier.CorrelationID)),
			zap.String("trace_id", obttrace.TraceIDFromContext(ctx)),
			zap.String("span_id", obttrace.SpanIDFromContext(ctx)),
			zap.String("grpc_service", grpcSvc),
			zap.String("grpc_method", grpcMethod),
			zap.String("grpc_code", codeStr),
			zap.Duration("latency", latency),
			zap.String("service", svcName),
		}
		if err != nil {
			cfg.Logger.Error("grpc unary request", append(fields, zap.Error(err))...)
		} else {
			cfg.Logger.Info("grpc unary request", fields...)
		}

		// ── Inject response metadata ─────────────────────────────────────────
		trailer := metadata.Pairs(
			ids.MetaResponseID, string(responseID),
			ids.MetaRequestID, string(carrier.RequestID),
		)
		_ = grpc.SetTrailer(ctx, trailer)

		return resp, err
	}
}

// ─── Streaming Server Interceptor ────────────────────────────────────────────

// StreamServerInterceptor returns a gRPC streaming server interceptor with
// the same observability capabilities as the unary interceptor.
func StreamServerInterceptor(cfg InterceptorConfig) grpc.StreamServerInterceptor {
	skipSet := makeSkipSet(cfg.SkipMethods)

	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if _, skip := skipSet[info.FullMethod]; skip {
			return handler(srv, ss)
		}

		ctx := propagation.GRPCExtract(ss.Context())
		carrier := ids.CarrierFromContext(ctx)

		grpcSvc, grpcMethod := splitMethod(info.FullMethod)
		ctx, span := cfg.Tracer.Start(ctx, info.FullMethod)
		span.SetGRPC(grpcSvc, grpcMethod, 0)

		if cfg.Metrics != nil {
			cfg.Metrics.InFlight.Inc()
			defer cfg.Metrics.InFlight.Dec()
		}

		start := time.Now()
		wrapped := &wrappedServerStream{ServerStream: ss, ctx: ctx}
		err := handler(srv, wrapped)
		latency := time.Since(start)
		grpcCode := statusCode(err)

		span.SetGRPC(grpcSvc, grpcMethod, uint32(grpcCode))
		if err != nil {
			span.RecordError(err)
		}
		span.End()

		svcName := cfg.Tracer.Config().ServiceName
		if cfg.Metrics != nil {
			lbls := []string{grpcSvc, grpcMethod, grpcCode.String(), svcName}
			cfg.Metrics.RequestsTotal.WithLabelValues(lbls...).Inc()
			cfg.Metrics.RequestDuration.WithLabelValues(lbls...).Observe(latency.Seconds())
		}

		cfg.Logger.Info("grpc stream request",
			zap.String("request_id", string(carrier.RequestID)),
			zap.String("trace_id", obttrace.TraceIDFromContext(ctx)),
			zap.String("grpc_service", grpcSvc),
			zap.String("grpc_method", grpcMethod),
			zap.String("grpc_code", grpcCode.String()),
			zap.Duration("latency", latency),
			zap.Bool("client_stream", info.IsClientStream),
			zap.Bool("server_stream", info.IsServerStream),
			zap.Error(err),
		)

		return err
	}
}

// ─── Client Interceptors ──────────────────────────────────────────────────────

// UnaryClientInterceptor propagates obslib IDs into outgoing gRPC calls
// and records client-side span + metrics.
func UnaryClientInterceptor(cfg InterceptorConfig) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = propagation.GRPCInject(ctx)

		grpcSvc, grpcMethod := splitMethod(method)
		ctx, span := cfg.Tracer.Start(ctx, method)
		defer func() { span.End() }()

		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		latency := time.Since(start)
		grpcCode := statusCode(err)

		span.SetGRPC(grpcSvc, grpcMethod, uint32(grpcCode))
		if err != nil {
			span.RecordError(err)
		}

		svcName := cfg.Tracer.Config().ServiceName
		if cfg.Metrics != nil {
			lbls := []string{grpcSvc, grpcMethod, grpcCode.String(), svcName}
			cfg.Metrics.RequestsTotal.WithLabelValues(lbls...).Inc()
			cfg.Metrics.RequestDuration.WithLabelValues(lbls...).Observe(latency.Seconds())
		}
		return err
	}
}

// StreamClientInterceptor propagates obslib IDs into outgoing gRPC streaming calls.
func StreamClientInterceptor(cfg InterceptorConfig) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		ctx = propagation.GRPCInject(ctx)
		_, span := cfg.Tracer.Start(ctx, method)
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			span.RecordError(err)
			span.End()
			return nil, err
		}
		return cs, nil
	}
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context { return w.ctx }

func splitMethod(full string) (service, method string) {
	// full = "/package.Service/Method"
	if len(full) == 0 || full[0] != '/' {
		return "unknown", full
	}
	s := full[1:]
	for i := range s {
		if s[i] == '/' {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}

func statusCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if s, ok := status.FromError(err); ok {
		return s.Code()
	}
	return codes.Unknown
}

func makeSkipSet(methods []string) map[string]struct{} {
	s := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		s[m] = struct{}{}
	}
	return s
}

// Ensure strconv is used (it's used in MetricsConfig indirectly).
var _ = strconv.Itoa
