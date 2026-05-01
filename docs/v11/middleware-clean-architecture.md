# V11 Middleware Clean Architecture

## Do

Use the canonical engine:

```go
reg := mwcore.EnterpriseRegistry(mwcore.EnterpriseOptions{
    Retry: retryCfg,
    Timeout: 2 * time.Second,
    RateLimit: rateLimitOptions,
})
handler := reg.Then(myHandler)
```

## Do not

Do not stack independent chains like this:

```text
middleware/http -> middleware/chain -> transport/core -> handler
```

Use one chain:

```text
middleware/core -> handler
```

## Ownership

| Concern | Owner |
|---|---|
| Retry, timeout, circuit breaker, backpressure, load shedding | `resilience/*` |
| Logging, metrics, tracing, correlation | `observability/*` |
| Ordering, dedupe, wrapper, adapter | `middleware/core` |
| HTTP/gRPC/NATS/Kafka protocol specifics | `transport/*` or `adapter/*` |

## Compatibility

Old packages such as `middleware/retry`, `middleware/timeout`, and `middleware/ratelimit` are retained as wrappers so existing services do not break. New code should compose them via `middleware/core`.
