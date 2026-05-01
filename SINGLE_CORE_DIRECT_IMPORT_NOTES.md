# Single Core Direct Import Clean Version

Canonical direct imports:

```go
import "github.com/PlatformCore/libpackage/observability/metrics"
import "github.com/PlatformCore/libpackage/resilience/retry"
import "github.com/PlatformCore/libpackage/core/errors"
import "github.com/PlatformCore/libpackage/resilience/ratelimit"
```

This version is built for direct package usage, not compat/wrapper usage.

Merged:
- `observability/gometrics` + `observability/metrics/basic` -> `observability/metrics`
- `resilience/retry/goretry` + `resilience/retry_backoff` -> `resilience/retry`
- `core/errors/goerror_compat` -> removed; use `core/errors`
- `ratelimit/limiter` + `ratelimit/enterprise` + selected legacy in-memory limiters -> root `ratelimit`

Removed duplicate package layers:
- `observability/gometrics`
- `observability/metrics/basic`
- `resilience/retry/goretry`
- `resilience/retry_backoff`
- `core/errors/goerror_compat`
- `ratelimit/compat`
- `ratelimit/legacy`
- `ratelimit/enterprise`
- `ratelimit/limiter`
- small duplicated ratelimit subpackages such as `policy`, `result`, `memory_store`, `redis_store`

Ratelimit naming:
- Core limiter: `ratelimit.New`, `ratelimit.Limiter`, `ratelimit.Policy`, `ratelimit.Key`, `ratelimit.Result`
- Enterprise distributed limiter: `ratelimit.EnterpriseLimiter`, `ratelimit.RedisLimiter`, `ratelimit.MemoryLimiter`, `ratelimit.CompositeLimiter`
- Legacy simple in-memory algorithms: `ratelimit.TokenBucketLimiter`, `ratelimit.FixedWindowLimiter`, `ratelimit.SlidingWindowLimiter`
- Legacy HTTP helpers that had name conflicts are prefixed with `Legacy...`
