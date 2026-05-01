# V11 Enterprise Cleanup Report

## Goal
V11 keeps the advanced V10 packages but removes duplicated responsibility in the middleware layer.

## What changed

### 1. Canonical middleware engine added
New package:

```text
middleware/core/
  canonical.go
  doc.go
  http_adapter.go
  pipeline.go
```

This is now the canonical middleware pipeline for cross-transport execution.

### 2. Duplicate responsibility guard
`middleware/core.Registry` registers middleware by normalized name and ignores duplicate names unless `Replace=true`.
This prevents accidental double application of:

- retry
- timeout
- ratelimit
- logging
- metrics
- tracing
- recovery

### 3. Clear responsibility rule

```text
middleware/core = order, dedupe, wrapper, adapter
resilience/*    = heavy retry/timeout/ratelimit/backpressure logic
observability/* = logging/metrics/tracing/correlation logic
transport/*     = protocol adapter only
```

### 4. Enterprise preset
Use `middleware/core.EnterpriseRegistry()` or `middleware/core.Enterprise()` to build the standard stack.

### 5. Removed stale backup artifacts
Removed non-Go backup/incoming files and temp files:

- `*.legacy.txt`
- `*_incoming.go.txt`
- `test.tmp`
- `test2.tmp`

No production Go implementation was removed.

## Canonical stack order

```text
Ingress       -> recovery
Policy        -> ratelimit
Resilience    -> retry, timeout
Observability -> logging, metrics, tracing
Business      -> handler
Egress        -> response/finalizer
```

## Files intentionally kept
Existing middleware packages remain for backward compatibility. They should now be used through `middleware/core` presets where possible.

## Remaining intentional duplicates
`VERSION` files may have identical content across modules. These are metadata files, not duplicate business logic.
