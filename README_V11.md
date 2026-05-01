# libpack v11 canonical enterprise cleanup

V11 is an upgrade from V10 focused on removing duplicated middleware responsibility while keeping the advanced enterprise packages intact.

## Main upgrade

New canonical middleware engine:

```text
middleware/core
```

It provides:

- deterministic ordered pipeline stages
- duplicate-name guard for retry/timeout/ratelimit/logging/metrics/tracing/recovery
- enterprise preset builder
- HTTP adapter backed by the same transport/core chain
- backward-compatible use of existing advanced middleware wrappers

## Recommended usage

```go
package main

import (
    "time"

    mwcore "github.com/PlatformCore/libpackage/middleware/core"
    retry "github.com/PlatformCore/libpackage/middleware/retry"
)

func buildStack() mwcore.Registry {
    cfg := retry.DefaultConfig()
    reg := mwcore.EnterpriseRegistry(mwcore.EnterpriseOptions{
        Retry: &cfg,
        Timeout: 2 * time.Second,
    })
    return *reg
}
```

## Rule

Do not make new retry/timeout/ratelimit logic inside middleware packages. Put heavy logic in `resilience/*` and wrap it through `middleware/core`.
