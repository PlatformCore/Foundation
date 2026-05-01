// Package telemetry bootstraps the full OpenTelemetry SDK for an enterprise
// service: tracer provider, metric provider, Prometheus exporter, and OTLP
// gRPC exporter. Call telemetry.Bootstrap() in main() and defer Shutdown().
package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	prometheusexporter "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config configures the telemetry bootstrap.
type Config struct {
	// ServiceName is the service name reported in traces and metrics. Required.
	ServiceName string
	// ServiceVersion is the service version (semver). Default: "0.0.0".
	ServiceVersion string
	// Environment is the deployment environment (production, staging, dev).
	Environment string
	// Region is the cloud region (us-east-1, eu-west-1, etc.).
	Region string

	// OTLPEndpoint is the gRPC endpoint of the OTLP collector (e.g. "localhost:4317").
	// If empty, trace export is disabled (useful in dev).
	OTLPEndpoint string
	// OTLPInsecure disables TLS for the OTLP connection. Default: false.
	OTLPInsecure bool

	// SamplingRatio is the trace sampling ratio (0.0 to 1.0). Default: 1.0.
	SamplingRatio float64

	// PrometheusRegisterer for metric export. Default: prometheus.DefaultRegisterer.
	PrometheusRegisterer prometheus.Registerer

	// BatchTimeout for the OTLP trace exporter. Default: 5s.
	BatchTimeout time.Duration

	// Logger for telemetry bootstrap events.
	Logger *zap.Logger
}

// Provider bundles all telemetry providers.
type Provider struct {
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	cfg            Config
	shutdown       []func(context.Context) error
}

// Bootstrap initialises the global OTel tracer and meter providers.
// It returns a Provider whose Shutdown() MUST be deferred in main().
func Bootstrap(ctx context.Context, cfg Config) (*Provider, error) {
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = "0.0.0"
	}
	if cfg.SamplingRatio == 0 {
		cfg.SamplingRatio = 1.0
	}
	if cfg.BatchTimeout == 0 {
		cfg.BatchTimeout = 5 * time.Second
	}
	if cfg.Logger == nil {
		cfg.Logger = zap.NewNop()
	}
	if cfg.PrometheusRegisterer == nil {
		cfg.PrometheusRegisterer = prometheus.DefaultRegisterer
	}

	p := &Provider{cfg: cfg}

	// ── Resource ─────────────────────────────────────────────────────────────
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
		resource.WithProcess(),
		resource.WithHost(),
		resource.WithOS(),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry: create resource: %w", err)
	}

	// ── Tracer Provider ───────────────────────────────────────────────────────
	tpOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(cfg.SamplingRatio),
		)),
	}

	if cfg.OTLPEndpoint != "" {
		exporterOpts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
			otlptracegrpc.WithTimeout(cfg.BatchTimeout),
		}
		if cfg.OTLPInsecure {
			exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
		}
		exporter, err := otlptracegrpc.New(ctx, exporterOpts...)
		if err != nil {
			return nil, fmt.Errorf("telemetry: create OTLP exporter: %w", err)
		}
		tpOpts = append(tpOpts,
			sdktrace.WithBatcher(exporter,
				sdktrace.WithBatchTimeout(cfg.BatchTimeout),
			),
		)
		p.shutdown = append(p.shutdown, exporter.Shutdown)
		cfg.Logger.Info("OTLP trace exporter configured", zap.String("endpoint", cfg.OTLPEndpoint))
	} else {
		cfg.Logger.Warn("OTLP endpoint not configured — traces will not be exported")
	}

	tp := sdktrace.NewTracerProvider(tpOpts...)
	otel.SetTracerProvider(tp)
	p.TracerProvider = tp
	p.shutdown = append(p.shutdown, tp.Shutdown)

	// ── Meter Provider (Prometheus exporter) ──────────────────────────────────
	promExporter, err := prometheusexporter.New(
		prometheusexporter.WithRegisterer(cfg.PrometheusRegisterer),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry: create prometheus exporter: %w", err)
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(promExporter),
	)
	otel.SetMeterProvider(mp)
	p.MeterProvider = mp
	p.shutdown = append(p.shutdown, mp.Shutdown)

	// ── Global Propagator ─────────────────────────────────────────────────────
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // W3C
		propagation.Baggage{},
	))

	cfg.Logger.Info("telemetry bootstrapped",
		zap.String("service", cfg.ServiceName),
		zap.String("version", cfg.ServiceVersion),
		zap.String("env", cfg.Environment),
		zap.Float64("sampling_ratio", cfg.SamplingRatio),
	)
	return p, nil
}

// Shutdown flushes and shuts down all telemetry providers.
// Call with a deadline context: ctx, cancel := context.WithTimeout(ctx, 10*time.Second).
func (p *Provider) Shutdown(ctx context.Context) error {
	var errs []error
	for i := len(p.shutdown) - 1; i >= 0; i-- {
		if err := p.shutdown[i](ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("telemetry shutdown errors: %v", errs)
	}
	p.cfg.Logger.Info("telemetry shut down cleanly")
	return nil
}

// Tracer returns a named tracer from the provider.
func (p *Provider) Tracer(name string) trace.Tracer {
	return p.TracerProvider.Tracer(name)
}

// Meter returns a named meter from the provider.
func (p *Provider) Meter(name string) metric.Meter {
	return p.MeterProvider.Meter(name)
}
