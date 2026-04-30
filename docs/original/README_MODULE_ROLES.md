# Libpackage Module Roles

à¹€à¸­à¸à¸ªà¸²à¸£à¸™à¸µà¹‰à¸ªà¸£à¸¸à¸›à¸«à¸™à¹‰à¸²à¸—à¸µà¹ˆà¸‚à¸­à¸‡à¸—à¸¸à¸à¹‚à¸¡à¸”à¸¹à¸¥à¸—à¸µà¹ˆà¸žà¸šà¸ˆà¸²à¸à¹„à¸Ÿà¸¥à¹Œ `go.mod` à¹ƒà¸•à¹‰ `libpackage` à¹à¸¥à¸°à¹€à¸Šà¹‡à¸à¸šà¸—à¸šà¸²à¸—à¸‹à¹‰à¸³ (overlap) à¹ƒà¸™à¸£à¸°à¸”à¸±à¸šà¸ªà¸–à¸²à¸›à¸±à¸•à¸¢à¸à¸£à¸£à¸¡

## 1) Foundation / Shared
- `github.com/driftappdev`: root aggregate module
- `github.com/PlatformCore/libpackage/core`: primitive à¸à¸¥à¸²à¸‡ (errors/logger/types/context/utils/result)
- `github.com/PlatformCore/libpackage/core/contracts`: à¸ªà¸±à¸à¸à¸² DTO/response/pagination/versioning
- `github.com/PlatformCore/libpackage/config`: config loading / env / defaults
- `github.com/PlatformCore/libpackage/core/result_legacy`: generic result envelope
- `github.com/PlatformCore/libpackage/runtime`: lifecycle/health/shutdown/runtime helpers
- `github.com/driftappdev/foundation/validator`: validation facade
- `github.com/PlatformCore/libpackage/clients`: client abstractions

## 2) Dependency Injection
- `github.com/PlatformCore/libpackage/orchestration/di`: DI umbrella
- `github.com/PlatformCore/libpackage/orchestration/di/container`: container
- `github.com/PlatformCore/libpackage/orchestration/di/module`: module registration
- `github.com/PlatformCore/libpackage/orchestration/di/provider`: provider wiring
- `github.com/PlatformCore/libpackage/orchestration/di/registry`: DI registry
- `github.com/PlatformCore/libpackage/orchestration/di/scope`: scope/lifetime

## 3) Messaging
- `github.com/PlatformCore/libpackage/messaging/inbox`: consumer-side intake queue
- `github.com/PlatformCore/libpackage/messaging/outbox`: producer-side outbox dispatch
- `github.com/PlatformCore/libpackage/messaging/dlq`: dead-letter queue
- `github.com/PlatformCore/libpackage/messaging/redrive`: DLQ redrive/replay to target
- `github.com/PlatformCore/libpackage/messaging/replay`: replay orchestration
- `github.com/PlatformCore/libpackage/messaging/idempotency`: idempotency guard
- `github.com/PlatformCore/libpackage/messaging/example_integration`: adapter/integration example
- `github.com/PlatformCore/libpackage/messaging/audit`: messaging audit events

## 4) Security
- `github.com/PlatformCore/libpackage/security`: security umbrella
- `github.com/PlatformCore/libpackage/security/jwt`: JWT
- `github.com/PlatformCore/libpackage/security/oauth2`: OAuth2
- `github.com/PlatformCore/libpackage/security/hash`: hashing
- `github.com/PlatformCore/libpackage/security/encryption`: encryption
- `github.com/PlatformCore/libpackage/security/permission`: permission checks
- `github.com/PlatformCore/libpackage/security/policy`: policy checks
- `github.com/PlatformCore/libpackage/security/threatdefense`: threat defense rules

## 5) Persistence
- `github.com/driftappdev/persistence`: persistence umbrella
- `github.com/PlatformCore/libpackage/persistence/tx`: transaction helpers
- `github.com/PlatformCore/libpackage/persistence/uow`: unit-of-work
- `github.com/PlatformCore/libpackage/persistence/distlock`: distributed lock

## 6) Observability / Telemetry
- `github.com/driftappdev/observability`: observability umbrella
- `github.com/PlatformCore/libpackage/observability/correlation`: correlation helpers
- `github.com/PlatformCore/libpackage/observability/healthcheck`: liveness/readiness/health
- `github.com/PlatformCore/libpackage/observability/span`: span helpers
- `github.com/PlatformCore/libpackage/observability/tracing`: tracing helpers
- `github.com/PlatformCore/libpackage/observability/performance`: performance metrics helpers
- `github.com/PlatformCore/libpackage/observability/profiler/sentinel`: profiler integration
- `github.com/PlatformCore/libpackage/observability/audit`: observability-side audit
- `github.com/driftappdev/observability/telemetry`: telemetry umbrella (OpenTelemetry-centric)

## 7) Platform / Plugins / Eventing
- `github.com/driftappdev/platform`: platform umbrella
- `github.com/PlatformCore/libpackage/platform/servicemesh`: service mesh integration
- `github.com/PlatformCore/libpackage/plugins`: plugin umbrella
- `github.com/PlatformCore/libpackage/plugins/common`: shared plugin contracts/util
- `github.com/PlatformCore/libpackage/plugins/engine`: plugin engine/runtime
- `github.com/PlatformCore/libpackage/platform/eventbus`: event bus abstractions

## 8) Resilience / Rate limiting
- `github.com/PlatformCore/libpackage/resilience`: resilience umbrella
- `github.com/PlatformCore/libpackage/resilience/retry`: retry
- `github.com/PlatformCore/libpackage/resilience/sanitizer`: sanitize/cleanup
- `github.com/PlatformCore/libpackage/resilience/validate`: validation utilities
- `github.com/PlatformCore/libpackage/resilience/validator`: validator helpers
- `github.com/PlatformCore/libpackage/resilience/pagination`: pagination helpers
- `github.com/PlatformCore/libpackage/resilience/cache`: cache resilience helpers
- `github.com/PlatformCore/libpackage/resilience/circuit`: circuit breaker
- `github.com/PlatformCore/libpackage/ratelimit`: rate limit umbrella
- `github.com/PlatformCore/libpackage/ratelimit/memory_store`: memory store for rate limit
- `github.com/PlatformCore/libpackage/ratelimit/enterprise`: enterprise rate-limit features

## 9) Middleware
- `github.com/driftappdev/middleware`: middleware umbrella
- `github.com/driftappdev/middleware/requestid`: request id middleware
- `github.com/driftappdev/middleware/timeout`: timeout middleware
- `github.com/driftappdev/middleware/logmid/logging-middleware`: logging middleware variant
- `github.com/driftappdev/middleware/adminshield/admin-middleware`: admin shield middleware variant
- `github.com/PlatformCore/libpackage/observability/logging-middleware`: standalone logging middleware
- `github.com/driftappdev/adminshield`: standalone adminshield (namespace à¸žà¸´à¸¡à¸žà¹Œà¸•à¹ˆà¸²à¸‡à¸ˆà¸²à¸ driftappdev)
- `github.com/driftappdev/auth`: standalone auth (namespace à¸žà¸´à¸¡à¸žà¹Œà¸•à¹ˆà¸²à¸‡à¸ˆà¸²à¸ driftappdev)

## 10) Legacy / Compatibility Facades (`go*`)
- `github.com/PlatformCore/libpackage/security/auth_middleware/goauth`
- `github.com/PlatformCore/libpackage/resilience/circuitbreaker/gocircuit`
- `github.com/PlatformCore/libpackage/core/errors/goerror_compat`
- `github.com/PlatformCore/libpackage/observability/logging/gologger`
- `github.com/PlatformCore/libpackage/observability/gometrics`
- `github.com/PlatformCore/libpackage/ratelimit/compat/goratelimit`
- `github.com/PlatformCore/libpackage/resilience/retry/goretry`
- `github.com/PlatformCore/libpackage/security/sanitizer/gosanitizer`
- `github.com/PlatformCore/libpackage/resilience/timeout/gotimeout`
- `github.com/PlatformCore/libpackage/observability/tracing/gotracing`

## 11) Folder Alias Map (à¸Šà¸·à¹ˆà¸­à¹‚à¸Ÿà¸¥à¹€à¸”à¸­à¸£à¹Œ != module path)
à¸à¸¥à¸¸à¹ˆà¸¡à¸™à¸µà¹‰à¹„à¸¡à¹ˆà¹ƒà¸Šà¹ˆà¹‚à¸¡à¸”à¸¹à¸¥à¹€à¸žà¸´à¹ˆà¸¡ à¹à¸•à¹ˆà¹€à¸›à¹‡à¸™à¹‚à¸Ÿà¸¥à¹€à¸”à¸­à¸£à¹Œà¸—à¸µà¹ˆà¸Šà¸µà¹‰à¹„à¸›à¸¢à¸±à¸‡ module path à¹€à¸Šà¸´à¸‡à¹‚à¸”à¹€à¸¡à¸™:
- `audit` -> `github.com/PlatformCore/libpackage/observability/audit`
- `auditv.1` -> `github.com/PlatformCore/libpackage/messaging/audit`
- `cache` -> `github.com/PlatformCore/libpackage/resilience/cache`
- `circuit` -> `github.com/PlatformCore/libpackage/resilience/circuit`
- `retry` -> `github.com/PlatformCore/libpackage/resilience/retry`
- `sanitizer` -> `github.com/PlatformCore/libpackage/resilience/sanitizer`
- `pagination` -> `github.com/PlatformCore/libpackage/resilience/pagination`
- `dlq` -> `github.com/PlatformCore/libpackage/messaging/dlq`
- `inbox` -> `github.com/PlatformCore/libpackage/messaging/inbox`
- `outbox` -> `github.com/PlatformCore/libpackage/messaging/outbox`
- `redrive` -> `github.com/PlatformCore/libpackage/messaging/redrive`
- `replay` -> `github.com/PlatformCore/libpackage/messaging/replay`
- `jwt` -> `github.com/PlatformCore/libpackage/security/jwt`
- `oauth2` -> `github.com/PlatformCore/libpackage/security/oauth2`
- `hash` -> `github.com/PlatformCore/libpackage/security/hash`
- `encryption` -> `github.com/PlatformCore/libpackage/security/encryption`
- `permission` -> `github.com/PlatformCore/libpackage/security/permission`
- `policy` -> `github.com/PlatformCore/libpackage/security/policy`
- `threatdefense` -> `github.com/PlatformCore/libpackage/security/threatdefense`
- `distlock` -> `github.com/PlatformCore/libpackage/persistence/distlock`
- `tx` -> `github.com/PlatformCore/libpackage/persistence/tx`
- `uow` -> `github.com/PlatformCore/libpackage/persistence/uow`
- `servicemesh` -> `github.com/PlatformCore/libpackage/platform/servicemesh`
- `common` -> `github.com/PlatformCore/libpackage/plugins/common`
- `engine` -> `github.com/PlatformCore/libpackage/plugins/engine`
- `enterprise` -> `github.com/PlatformCore/libpackage/ratelimit/enterprise`
- `performance` -> `github.com/PlatformCore/libpackage/observability/performance`
- `sentinelprofiler` -> `github.com/PlatformCore/libpackage/observability/profiler/sentinel`
- `ratelimitX` -> `github.com/PlatformCore/libpackage/ratelimit`

---

## Overlap Check (à¸šà¸—à¸šà¸²à¸—à¸‹à¹‰à¸³)

### A) à¸‹à¹‰à¸³à¹à¸šà¸š â€œAlias Moduleâ€
à¸à¸¥à¸¸à¹ˆà¸¡à¸™à¸µà¹‰à¸‹à¹‰à¸³à¸šà¸—à¸šà¸²à¸—à¹‚à¸”à¸¢à¸•à¸±à¹‰à¸‡à¹ƒà¸ˆ (à¸Šà¸·à¹ˆà¸­à¸ªà¸±à¹‰à¸™ vs à¸Šà¸·à¹ˆà¸­à¹‚à¸”à¹€à¸¡à¸™à¹€à¸•à¹‡à¸¡):
- `inbox/outbox/dlq/redrive/replay` vs `messaging/*`
- `retry/sanitizer/cache/circuit/pagination` vs `resilience/*`
- `jwt/oauth2/hash/encryption/permission/policy/threatdefense` vs `security/*`
- `servicemesh` vs `platform/servicemesh`
- `distlock/tx/uow` vs `persistence/*`
- `audit/performance/sentinelprofiler` vs `observability/*`
- `ratelimitX` vs `ratelimit`

### B) à¸‹à¹‰à¸³à¹€à¸Šà¸´à¸‡à¹‚à¸”à¹€à¸¡à¸™ (à¹„à¸¡à¹ˆà¸ˆà¸³à¹€à¸›à¹‡à¸™à¸•à¹‰à¸­à¸‡à¸¡à¸µà¸—à¸±à¹‰à¸‡à¸„à¸¹à¹ˆ)
- `telemetry` vs `observability/*` (tracing/metrics overlap à¸ªà¸¹à¸‡)
- `plugins/*` vs `platform/*` (loader/registry/provider/hooks overlap)
- `validator` vs `resilience/validate` vs `resilience/validator`
- `logging-middleware` vs `middleware/logmid/logging-middleware`
- `auditv.1 (messaging/audit)` vs `audit (observability/audit)` à¸•à¹‰à¸­à¸‡à¹à¸¢à¸à¸‚à¸­à¸šà¹€à¸‚à¸•à¹ƒà¸«à¹‰à¸Šà¸±à¸”

### C) à¸ˆà¸¸à¸”à¹€à¸ªà¸µà¹ˆà¸¢à¸‡à¸„à¸§à¸²à¸¡à¹„à¸¡à¹ˆà¸ªà¸¡à¹ˆà¸³à¹€à¸ªà¸¡à¸­
- namespace à¸ªà¸°à¸à¸”à¸•à¹ˆà¸²à¸‡à¸à¸±à¸™: `github.com/diftappdev/...` vs `github.com/driftappdev/...`
- à¸šà¸²à¸‡à¹‚à¸¡à¸”à¸¹à¸¥à¸¢à¸±à¸‡à¹€à¸›à¹‡à¸™ legacy alias à¸ˆà¸³à¸™à¸§à¸™à¸¡à¸²à¸ à¸—à¸³à¹ƒà¸«à¹‰à¸”à¸¹à¹€à¸«à¸¡à¸·à¸­à¸™à¸‹à¹‰à¸³à¹€à¸¢à¸­à¸°à¸à¸§à¹ˆà¸²à¸„à¸§à¸²à¸¡à¸ˆà¸£à¸´à¸‡

## Recommended Direction
- à¸–à¹‰à¸²à¸•à¹‰à¸­à¸‡à¸à¸²à¸£à¸¥à¸”à¸‹à¹‰à¸³: à¹€à¸¥à¸·à¸­à¸ â€œcanonical pathâ€ à¸•à¹ˆà¸­à¹‚à¸”à¹€à¸¡à¸™à¸¥à¸° 1 à¸Šà¸¸à¸” (à¹€à¸Šà¹ˆà¸™ `messaging/*`, `security/*`, `resilience/*`, `observability/*`, `platform/*`) à¹à¸¥à¹‰à¸§à¸„à¸‡ alias à¹€à¸‰à¸žà¸²à¸°à¸—à¸µà¹ˆà¸ˆà¸³à¹€à¸›à¹‡à¸™à¸•à¹ˆà¸­ backward compatibility
- à¸à¸³à¸«à¸™à¸” policy à¸Šà¸±à¸”à¹€à¸ˆà¸™: module à¹ƒà¸«à¸¡à¹ˆà¸•à¹‰à¸­à¸‡à¹€à¸‚à¹‰à¸²à¸Šà¸¸à¸” canonical à¸à¹ˆà¸­à¸™à¹€à¸ªà¸¡à¸­
- à¹à¸à¹‰ namespace à¹ƒà¸«à¹‰à¹€à¸«à¸¥à¸·à¸­ `driftappdev` à¹à¸šà¸šà¹€à¸”à¸µà¸¢à¸§à¸—à¸±à¹‰à¸‡ repo




