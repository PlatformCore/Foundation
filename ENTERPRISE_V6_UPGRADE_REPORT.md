# libpack v6 enterprise ultra plus - single core direct import

## Goal
Keep direct package imports, keep existing capabilities, and add missing enterprise-grade building blocks without restoring compat/wrapper/legacy duplication.

## Added packages / files

### cache
- `cache/cache.go`
- `cache/memory.go`
- `cache/readthrough.go`

Adds memory LRU cache, TTL, cleanup, single-flight read-through loading, basic cache metrics.

### config
- `config/validator/validator.go`
- `config/override/override.go`

Adds struct-tag validation and environment override support.

### observability
- `observability/healthcheck/dependency.go`
- `observability/slo/slo.go`

Adds dependency health checks and SLO/error-budget tracking.

### coordination
- `coordination/lease.go`
- `coordination/leader_election.go`

Adds lease/fencing-token based distributed coordination interfaces and in-memory implementation.

### resilience
- `resilience/adaptive_limiter/adaptive_limiter.go`

Adds adaptive concurrency limiter.

### security
- `security/secrets/secrets.go`

Adds secret store interface, memory store, and AES-GCM envelope encryption helpers.

### messaging
- `messaging/poison/poison.go`

Adds poison-message classification and quarantine decision support.

### runtime
- `runtime/supervisor/supervisor.go`

Adds restart supervisor for long-running workers.

## Design rules preserved
- No compat layer restored.
- No wrapper-first package paths.
- Heavy logic remains in the direct package where users import it.
- Existing packages are kept.
- New modules are separated only when they are real top-level capabilities (`cache`, `coordination`), not compatibility shims.
