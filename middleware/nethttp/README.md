# enterprise/middleware

Enterprise-grade, framework-agnostic HTTP middleware package for Go services.  
Compatible with `net/http` and any router that accepts `http.Handler` (Chi, Gorilla, Echo adapters, etc.).

## Features

| File | Middleware | Description |
|------|-----------|-------------|
| `requestid.go` | `WithRequestID` | Injects X-Request-ID into context & response |
| `logger.go` | `WithLogger` | Structured zap logging (method, path, status, latency) |
| `recovery.go` | `WithRecovery` | Panic recovery → 500 JSON response |
| `security.go` | `WithSecurity` | OWASP security headers (HSTS, CSP, X-Frame-Options…) |
| `cors.go` | `WithCORS` | Configurable CORS with preflight handling |
| `compress.go` | `WithCompress` | Gzip response compression with content-type filtering |
| `tracing.go` | `WithTracing` | OpenTelemetry distributed tracing (W3C TraceContext) |
| `metrics.go` | `WithMetrics` | Prometheus: request count, latency histogram, in-flight |
| `auth.go` | `WithAuth` / `RequireRoles` | JWT authentication (HS256/RS256) + RBAC |
| `ratelimit.go` | `WithRateLimit` | Token-bucket rate limiting per IP / user / tenant |
| `circuit_breaker.go` | `WithCircuitBreaker` | Closed → Open → Half-Open circuit breaker |
| `timeout.go` | `WithTimeout` | Per-request context timeout |
| `cache.go` | `WithCache` | Pluggable response caching (in-memory or Redis) |
| `validator.go` | `WithBodyValidator[T]` | Generic JSON body decode + validation |
| `health.go` | `WithHealth` | `/healthz` liveness + `/readyz` readiness probes |

## Installation

```bash
go get github.com/PlatformCore/Foundation/middleware/nethttp
```

## Quick Start

```go
package main

import (
    "net/http"
    "time"

    "go.uber.org/zap"
    "github.com/PlatformCore/Foundation/middleware/nethttp"
)

func main() {
    logger, _ := zap.NewProduction()

    chain := middleware.New(
        middleware.WithRequestID(),
        middleware.WithLogger(logger),
        middleware.WithRecovery(middleware.RecoveryOptions{Logger: logger}),
        middleware.WithSecurity(),
        middleware.WithCORS(middleware.CORSConfig{AllowedOrigins: []string{"*"}}),
        middleware.WithCompress(),
        middleware.WithTracing(middleware.TracingConfig{ServiceName: "my-service"}),
        middleware.WithMetrics(middleware.MetricsConfig{}),
        middleware.WithAuth(middleware.AuthConfig{
            HMACSecret: []byte("my-secret"),
        }),
        middleware.WithRateLimit(middleware.RateLimitConfig{
            RequestsPerSecond: 100,
            Burst:             200,
        }),
        middleware.WithTimeout(middleware.TimeoutConfig{Timeout: 30 * time.Second}),
    )

    mux := http.NewServeMux()
    middleware.WithHealth(middleware.HealthConfig{}, mux)
    mux.Handle("/api/", chain.ThenFunc(myHandler))

    http.ListenAndServe(":8080", mux)
}
```

## Recommended Chain Order

```
RequestID → Logger → Recovery → Security → CORS → Compress →
Tracing → Metrics → Auth → RateLimit → CircuitBreaker → Timeout → Cache
```

> **Why this order?**  
> - `RequestID` must be first so all subsequent middleware can log it.  
> - `Logger` wraps everything to capture the final status code.  
> - `Recovery` sits inside Logger so panics are logged before recovery.  
> - `Auth` comes after observability so unauthenticated requests are still traced/logged.  
> - `RateLimit` comes after Auth so you can limit per-user.  
> - `Cache` is last so Auth/RateLimit still run on cache hits.

## Context Helpers

```go
// Get the request ID from any handler or middleware.
id := middleware.GetRequestID(r.Context())

// Get validated JWT claims.
claims := middleware.GetClaims(r.Context())
fmt.Println(claims.UserID, claims.Roles, claims.TenantID)

// Get a validated + decoded request body.
body, ok := middleware.GetValidatedBody[CreateUserRequest](r.Context())
```

## Generic Body Validation

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

func (r CreateUserRequest) Validate() error {
    var errs middleware.ValidationErrors
    if e := middleware.RequiredField("name", r.Name); e != nil {
        errs = append(errs, *e)
    }
    if e := middleware.IsEmail("email", r.Email); e != nil {
        errs = append(errs, *e)
    }
    if e := middleware.IsOneOf("role", r.Role, "admin", "user", "viewer"); e != nil {
        errs = append(errs, *e)
    }
    if len(errs) > 0 {
        return errs
    }
    return nil
}

// Register:
chain.Append(middleware.WithBodyValidator[CreateUserRequest](middleware.ValidatorConfig{}))
```

## Circuit Breaker States

```
Closed ──(MaxFailures exceeded)──► Open
Open   ──(OpenTimeout elapsed)───► Half-Open
Half-Open ──(probe succeeds)────► Closed
Half-Open ──(probe fails)───────► Open
```

## Cache Backends

The default backend is an in-process `sync.Map`-based LRU.  
For distributed caching, implement `middleware.CacheBackend`:

```go
type CacheBackend interface {
    Get(ctx context.Context, key string) ([]byte, bool)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

## License

MIT
