# V10 Canonical Single-Core Cleanup Report

This version removes compatibility/facade paths because the library has not yet been adopted by downstream services.
No feature set was intentionally removed; duplicate entry points were merged into the canonical packages.

## Canonical packages

| Old / duplicate path | V10 canonical path | Action |
|---|---|---|
| `core/logger` | `observability/logging` | merged old field/slog helpers into canonical logging and removed duplicate folder |
| `middleware/nethttp` | `middleware/http` | removed duplicate path and updated imports |
| root `ratelimit` | `resilience/ratelimit` | kept limiter algorithms/features in resilience package and removed duplicate root module |
| `plugins/engine` | `plugins/runtimeengine` | moved implementation to clearer runtimeengine package and removed generic engine path |
| `middleware/obshttp` | `middleware/http` | merged observability middleware into canonical HTTP middleware package |

## Important implementation notes

- HTTP rate limiting now delegates to `resilience/ratelimit`; it does not keep a separate limiter engine in middleware.
- HTTP observability now lives in `middleware/http/observability.go`.
- Plugin runtime engine implementation now lives in `plugins/runtimeengine/engine.go`.
- `go.work` no longer references the removed root `ratelimit` module.
- Stale compatibility folders were deleted rather than retained as deprecated wrappers.

## Scan result

- Removed duplicate directories still present: `[]`
- Old canonical import references found in Go files: `[]`
- Total files: `471`
- Go files: `360`
- Total lines: `54284`
- Go lines: `44252`
