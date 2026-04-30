# enterprise/obslib — Enterprise Observability Library

Production-grade, **zero-framework** shared library for Go microservices at Google/Netflix/Uber scale.  
Handles ID generation, cross-process propagation, distributed tracing, Prometheus metrics,
and structured logging across **HTTP**, **gRPC**, and **Event-driven** transports — all from a single package.

---

## Architecture

```
obslib/
├── pkg/
│   ├── ids/              # Request/Response/Trace/Span/Correlation/Session ID types & context helpers
│   ├── trace/            # Enterprise OTel Tracer wrapper — enriched spans with semantic attrs
│   ├── propagation/      # Cross-protocol ID propagation (HTTP ↔ gRPC ↔ Event)
│   ├── registry/         # Thread-safe component registry with health polling & lifecycle hooks
│   ├── telemetry/        # OTel SDK bootstrap (TracerProvider + MeterProvider + Prometheus)
│   └── middleware/
│       ├── http/         # net/http middleware chain (propagation + tracing + metrics + logging)
│       ├── grpc/         # gRPC unary + streaming interceptors (server + client)
│       └── event/        # Event/message middleware (Kafka, NATS, RabbitMQ, Pub/Sub, SQS)
├── internal/
│   ├── clock/            # Testable clock abstraction
│   └── pool/             # sync.Pool byte buffer helpers
└── example/              # Complete multi-transport wiring example
```

---

## Packages

### `pkg/ids` — ID Types & Context

| Type | Format | Purpose |
|------|--------|---------|
| `RequestID` | `<13hex-ms><6hex-seq>-<32hex-uuid>` | Unique per inbound request |
| `ResponseID` | `rsp_<13hex-ms>_<16hex-random>` | Unique per outbound response |
| `TraceID` | 16-byte / 32-hex (W3C) | Links all spans in a trace |
| `SpanID` | 8-byte / 16-hex (W3C) | Individual operation span |
| `CorrelationID` | `cor_<base64url>` | Groups related requests (saga/workflow) |
| `SessionID` | `ses_<base64url>` | User session |
| `TraceParent` | W3C `traceparent` format | Full W3C propagation |

```go
// Generate IDs
rid := ids.NewRequestID()
tid := ids.NewTraceID()
cid := ids.NewCorrelationID()

// Bulk context injection
carrier := ids.Carrier{RequestID: rid, TraceID: tid, CorrelationID: cid}
ctx = ids.WithCarrier(ctx, carrier)

// Extraction
carrier = ids.CarrierFromContext(ctx)
fmt.Println(carrier.RequestID, carrier.TraceID)
```

### `pkg/propagation` — Multi-Protocol Propagation

Supports 3 strategies simultaneously, with priority order:
1. **X-Obs-Context** (enterprise JSON envelope — single header, all IDs)
2. **W3C traceparent** (OpenTelemetry standard)
3. **B3 multi-header** (Zipkin/legacy interop)

```go
// HTTP — extract on server
ctx = propagation.HTTPExtract(r.Context(), r.Header)

// HTTP — inject on client
propagation.HTTPInject(ctx, req.Header)

// gRPC — extract on server
ctx = propagation.GRPCExtract(ctx)          // reads metadata.FromIncomingContext

// gRPC — inject on client
ctx = propagation.GRPCInject(ctx)            // adds to metadata.NewOutgoingContext

// Event/Message (Kafka, NATS, etc.)
ctx = propagation.EventExtract(ctx, msg.Headers())
propagation.EventInject(ctx, outMsg.Headers())

// Kafka adapter
headers := propagation.NewKafkaHeaders(rawHeaderMap)
ctx = propagation.EventExtract(ctx, headers)
```

### `pkg/trace` — Enterprise Span API

```go
tracer := trace.NewTracer("my-service", trace.SpanConfig{
    ServiceName: "order-service", Environment: "production", Region: "ap-southeast-1",
})

ctx, span := tracer.Start(ctx, "checkout.process")
span.
    SetUser("user-123", "tenant-abc").
    SetHTTP("POST", "/api/v1/orders", 201).
    SetDB("postgres", "INSERT INTO orders ...").
    AddEvent("order.created", attribute.String("order.id", "ord-xyz"))
span.End()

// Error recording with stack trace
ctx, span = tracer.Start(ctx, "payment.charge")
if err := chargeCard(ctx, amount); err != nil {
    span.EndWithError(err) // records error type, message, and full stack
}
```

### `pkg/registry` — Component Registry

```go
reg := registry.New(
    registry.WithLogger(logger),
    registry.WithHealthInterval(30 * time.Second),
    registry.WithLifecycleHook(func(event string, snap registry.ComponentSnapshot) {
        logger.Info("component "+event, zap.String("name", snap.Name))
    }),
)
defer reg.Shutdown()

reg.MustRegister(&registry.Component{
    Name:    "redis-cache",
    Kind:    registry.KindMetricCollector,
    Version: "7.2",
    HealthFn: func(ctx context.Context) error {
        return redisClient.Ping(ctx).Err()
    },
})

// Query by kind or tag
tracers := reg.GetByKind(registry.KindTracer)
fmt.Println(reg.IsHealthy()) // true if all health checks pass
```

### `pkg/telemetry` — OTel Bootstrap

```go
tel, err := telemetry.Bootstrap(ctx, telemetry.Config{
    ServiceName:    "order-service",
    ServiceVersion: "1.4.2",
    Environment:    "production",
    Region:         "ap-southeast-1",
    OTLPEndpoint:   "otel-collector:4317",  // leave empty to disable export (dev)
    OTLPInsecure:   false,
    SamplingRatio:  0.1,                    // 10% in production
    Logger:         logger,
})
defer tel.Shutdown(ctx) // flushes OTLP exporter + Prometheus
```

### HTTP Middleware Chain

```go
chain := obshttp.Chain(
    obshttp.WithPropagation(),                    // ID extract → context → response headers
    obshttp.WithRecovery(logger, tracer),          // Panic → 500 JSON + span error
    obshttp.WithObservability(obshttp.ObsConfig{   // Trace + metrics + structured log
        Tracer:        tracer,
        Logger:        logger,
        Metrics:       obshttp.DefaultMetrics(nil, "obs", "http"),
        SkipPaths:     []string{"/healthz", "/metrics"},
        SlowThreshold: 500 * time.Millisecond,
    }),
)

mux.Handle("/api/v1/orders", chain(http.HandlerFunc(handler)))
```

### gRPC Interceptors

```go
grpcSrv := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        obsgrpc.UnaryServerInterceptor(obsgrpc.InterceptorConfig{
            Tracer:      tracer,
            Logger:      logger,
            Metrics:     obsgrpc.DefaultGRPCMetrics(nil, "obs"),
            SkipMethods: []string{"/grpc.health.v1.Health/Check"},
        }),
    ),
    grpc.ChainStreamInterceptor(
        obsgrpc.StreamServerInterceptor(cfg),
    ),
)

// Client side — propagate context into outgoing calls
conn, _ := grpc.Dial(addr,
    grpc.WithChainUnaryInterceptor(obsgrpc.UnaryClientInterceptor(cfg)),
    grpc.WithChainStreamInterceptor(obsgrpc.StreamClientInterceptor(cfg)),
)
```

### Event Middleware (Kafka / NATS / RabbitMQ / SQS)

```go
handler := event.Chain(
    event.WithPropagation(),       // extract IDs + create child span ID
    event.WithRecovery(logger),    // panic → error (so retry fires)
    event.WithRetry(event.RetryConfig{
        MaxAttempts: 3,
        Backoff:     event.ExponentialBackoff(100 * time.Millisecond),
        OnDLQ: func(ctx context.Context, msg event.Message, err error) {
            // send to dead-letter queue
        },
    }),
    event.WithObservability(event.EventConfig{
        Tracer: tracer, Logger: logger,
        Metrics: event.DefaultEventMetrics(nil, "obs"),
    }),
)(myHandler)

// Publish with automatic context propagation
publisher := event.InstrumentedPublisher(tracer, logger, brokerPublish)
publisher(ctx, "orders.created", headers, payload)
```

---

## Recommended Middleware Order

### HTTP
```
WithPropagation → WithRecovery → WithObservability → [custom] → Handler
```

### gRPC
```
UnaryServerInterceptor (propagation + tracing + metrics + logging inside)
```

### Event
```
WithPropagation → WithRecovery → WithRetry → WithObservability → Handler
```

---

## Metrics Emitted

| Protocol | Metric | Labels |
|----------|--------|--------|
| HTTP | `obs_http_requests_total` | method, path, status, service |
| HTTP | `obs_http_request_duration_seconds` | method, path, status, service |
| HTTP | `obs_http_requests_in_flight` | — |
| HTTP | `obs_http_response_size_bytes` | method, path, service |
| gRPC | `obs_grpc_requests_total` | grpc_service, grpc_method, grpc_code, service |
| gRPC | `obs_grpc_request_duration_seconds` | grpc_service, grpc_method, grpc_code, service |
| gRPC | `obs_grpc_requests_in_flight` | — |
| Event | `obs_event_messages_processed_total` | topic, status, service |
| Event | `obs_event_processing_duration_seconds` | topic, status, service |
| Event | `obs_event_messages_in_flight` | — |
| Event | `obs_event_retries_total` | topic, service |
| Event | `obs_event_dlq_total` | topic, service |

---

## Installation

```bash
go get github.com/PlatformCore/libpackage/observability/obslib
```

## License

MIT
