# libpack enterprise ultra single-core merge report

This build keeps the direct-import layout and adds the advanced logic detected from the two supplied missing-parts ZIPs without restoring compat/legacy wrapper packages.

## Direct package policy

Use these packages directly from microservices:

```go
import "github.com/PlatformCore/libpackage/observability/metrics"
import "github.com/PlatformCore/libpackage/resilience/retry"
import "github.com/PlatformCore/libpackage/core/errors"
import "github.com/PlatformCore/libpackage/resilience/ratelimit"
import "github.com/PlatformCore/libpackage/middleware/propagation"
import "github.com/PlatformCore/libpackage/middleware/registry"
```

## Advanced functions merged into the current core packages

### ratelimit
- Advanced key shape: namespace, identity, route, method, tenant, dimensions.
- Key validation and deterministic normalized key formatting.
- Additional strategies: fixed window, sliding window, token bucket, leaky bucket constant.
- Policy fields for burst and refill rate.
- Result HTTP headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, `Retry-After`.
- EvalStore model with StoreRequest/StoreResponse.
- Direct in-package MemoryStore supporting fixed-window, sliding-window, token-bucket.
- Direct in-package RedisStore interface for Redis-compatible fixed/sliding enforcement.
- Limiter now automatically uses EvalStore when the selected store supports advanced strategy evaluation.

### retry
- Decorrelated jitter strategy added to the direct `resilience/retry` package.
- RetryCounter made atomic for concurrent use.
- Existing retry builder, typed retry, predicates, per-attempt timeout, and backoff modes remain in the same direct package.

### middleware registry
- Added direct component registry in `middleware/registry` for runtime component health/lifecycle tracking.
- No wrapper/compat path required.

### metrics/errors
- The direct metrics package already contained the heavy gometrics registry/counter/gauge/histogram/summary/timer/prometheus logic plus the clean simple registry API.
- The direct errors package already contained BaseError, category, stacktrace, API response mapping, helpers, and aggregate errors.

## What was intentionally not restored

- `compat/*` folders and old duplicate wrapper paths.
- Whole legacy files copied as-is.
- Historical `_module_backup` go.mod/go.sum files.

The intent is function-level preservation inside the clean package structure, not path-level restoration.
