// Command example demonstrates a complete obslib integration for a service
// that exposes HTTP, gRPC, and event consumer endpoints simultaneously.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/PlatformCore/Foundation/middleware/ids"
	"github.com/PlatformCore/Foundation/middleware/event"
	obshttp "github.com/PlatformCore/Foundation/middleware/obshttp"
	obsgrpc "github.com/PlatformCore/Foundation/middleware/obsgrpc"
	"github.com/PlatformCore/Foundation/middleware/propagation"
	"github.com/PlatformCore/Foundation/middleware/registry"
	"github.com/PlatformCore/Foundation/observability/telemetry/provider"
	obttrace "github.com/PlatformCore/Foundation/middleware/trace"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// ── Telemetry Bootstrap ──────────────────────────────────────────────────
	tel, err := telemetry.Bootstrap(ctx, telemetry.Config{
		ServiceName:    "order-service",
		ServiceVersion: "1.4.2",
		Environment:    "production",
		Region:         "ap-southeast-1",
		OTLPEndpoint:   "otel-collector:4317",
		OTLPInsecure:   false,
		SamplingRatio:  0.1,
		Logger:         logger,
	})
	if err != nil {
		log.Fatalf("telemetry bootstrap: %v", err)
	}
	defer func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = tel.Shutdown(shutCtx)
	}()

	// ── Enterprise Tracer ────────────────────────────────────────────────────
	tracer := obttrace.NewTracer("order-service", obttrace.SpanConfig{
		ServiceName:    "order-service",
		ServiceVersion: "1.4.2",
		Environment:    "production",
		Region:         "ap-southeast-1",
	})

	// ── Component Registry ───────────────────────────────────────────────────
	reg := registry.New(
		registry.WithLogger(logger),
		registry.WithHealthInterval(30*time.Second),
		registry.WithLifecycleHook(func(evt string, snap registry.ComponentSnapshot) {
			logger.Info("registry event",
				zap.String("event", evt),
				zap.String("component", snap.Name),
			)
		}),
	)
	defer reg.Shutdown()
	reg.MustRegister(&registry.Component{
		Name: "otel-tracer", Kind: registry.KindTracer,
		Description: "OpenTelemetry enterprise tracer", Version: "1.24.0",
		Value: tracer,
	})

	// ── HTTP Server ──────────────────────────────────────────────────────────
	httpMetrics := obshttp.DefaultMetrics(nil, "obs", "http")
	httpChain := obshttp.Chain(
		obshttp.WithPropagation(),
		obshttp.WithRecovery(logger, tracer),
		obshttp.WithObservability(obshttp.ObsConfig{
			Tracer:            tracer,
			Logger:            logger,
			Metrics:           httpMetrics,
			SkipPaths:         []string{"/healthz", "/readyz", "/metrics"},
			SlowThreshold:     500 * time.Millisecond,
			LogRequestHeaders: []string{"X-Tenant-ID", "X-User-ID"},
		}),
	)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"up"}`))
	})
	mux.Handle("/api/v1/orders", httpChain(http.HandlerFunc(ordersHandler)))

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		logger.Info("HTTP server listening", zap.String("addr", ":8080"))
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// ── gRPC Server ──────────────────────────────────────────────────────────
	grpcMetrics := obsgrpc.DefaultGRPCMetrics(nil, "obs")
	interceptorCfg := obsgrpc.InterceptorConfig{
		Tracer:      tracer,
		Logger:      logger,
		Metrics:     grpcMetrics,
		SkipMethods: []string{"/grpc.health.v1.Health/Check"},
	}
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(obsgrpc.UnaryServerInterceptor(interceptorCfg)),
		grpc.ChainStreamInterceptor(obsgrpc.StreamServerInterceptor(interceptorCfg)),
	)
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("gRPC listen: %v", err)
	}
	go func() {
		logger.Info("gRPC server listening", zap.String("addr", ":9090"))
		if err := grpcSrv.Serve(lis); err != nil {
			logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	// ── Event Consumer ───────────────────────────────────────────────────────
	eventMetrics := event.DefaultEventMetrics(nil, "obs")
	eventChain := event.Chain(
		event.WithPropagation(),
		event.WithRecovery(logger),
		event.WithRetry(event.RetryConfig{
			MaxAttempts: 3,
			Backoff:     event.ExponentialBackoff(100 * time.Millisecond),
			Metrics:     eventMetrics,
			ServiceName: "order-service",
			OnDLQ: func(ctx context.Context, msg event.Message, err error) {
				carrier := ids.CarrierFromContext(ctx)
				logger.Error("message sent to DLQ",
					zap.String("topic", msg.Topic()),
					zap.String("request_id", string(carrier.RequestID)),
					zap.Error(err),
				)
			},
		}),
		event.WithObservability(event.EventConfig{
			Tracer: tracer, Logger: logger,
			Metrics: eventMetrics, ServiceName: "order-service",
		}),
	)
	_ = eventChain

	// Instrumented publisher — inject context into every outgoing message.
	publisher := event.InstrumentedPublisher(tracer, logger,
		func(ctx context.Context, topic string, headers propagation.MessageHeader, body []byte) error {
			logger.Debug("publishing", zap.String("topic", topic))
			return nil
		},
	)
	_ = publisher

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	<-ctx.Done()
	logger.Info("shutting down...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutCtx)
	grpcSrv.GracefulStop()
	logger.Info("shutdown complete")
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	carrier := ids.CarrierFromContext(ctx)
	fmt.Fprintf(w, `{"request_id":%q,"correlation_id":%q,"trace_id":%q,"span_id":%q}`,
		carrier.RequestID,
		carrier.CorrelationID,
		obttrace.TraceIDFromContext(ctx),
		obttrace.SpanIDFromContext(ctx),
	)
}
