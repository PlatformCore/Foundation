# Libpackage Module Roles

à¹€à¸­à¸à¸ªà¸²à¸£à¸™à¸µà¹‰à¸ªà¸£à¸¸à¸›à¸«à¸™à¹‰à¸²à¸—à¸µà¹ˆà¸‚à¸­à¸‡à¸—à¸¸à¸à¹‚à¸¡à¸”à¸¹à¸¥à¸—à¸µà¹ˆà¸žà¸šà¸ˆà¸²à¸à¹„à¸Ÿà¸¥à¹Œ `go.mod` à¹ƒà¸•à¹‰ `libpackage` à¹à¸¥à¸°à¹€à¸Šà¹‡à¸à¸šà¸—à¸šà¸²à¸—à¸‹à¹‰à¸³ (overlap) à¹ƒà¸™à¸£à¸°à¸”à¸±à¸šà¸ªà¸–à¸²à¸›à¸±à¸•à¸¢à¸à¸£à¸£à¸¡

## 1) Foundation / Shared
- `github.com/driftappdev`: root aggregate module
- `github.com/PlatformCore/Foundation/core`: primitive à¸à¸¥à¸²à¸‡ (errors/logger/types/context/utils/result)
- `github.com/PlatformCore/Foundation/core/contracts`: à¸ªà¸±à¸à¸à¸² DTO/response/pagination/versioning
- `github.com/PlatformCore/Foundation/config`: config loading / env / defaults
- `github.com/PlatformCore/Foundation/core/result_legacy`: generic result envelope
- `github.com/PlatformCore/Foundation/runtime`: lifecycle/health/shutdown/runtime helpers
- `github.com/driftappdev/foundation/validator`: validation facade
- `github.com/PlatformCore/Foundation/clients`: client abstractions

## 2) Dependency Injection
- `github.com/PlatformCore/Foundation/orchestration/di`: DI umbrella
- `github.com/PlatformCore/Foundation/orchestration/di/container`: container
- `github.com/PlatformCore/Foundation/orchestration/di/module`: module registration
- `github.com/PlatformCore/Foundation/orchestration/di/provider`: provider wiring
- `github.com/PlatformCore/Foundation/orchestration/di/registry`: DI registry
- `github.com/PlatformCore/Foundation/orchestration/di/scope`: scope/lifetime

## 3) Messaging
- `github.com/PlatformCore/Foundation/messaging/inbox`: consumer-side intake queue
- `github.com/PlatformCore/Foundation/messaging/outbox`: producer-side outbox dispatch
- `github.com/PlatformCore/Foundation/messaging/dlq`: dead-letter queue
- `github.com/PlatformCore/Foundation/messaging/redrive`: DLQ redrive/replay to target
- `github.com/PlatformCore/Foundation/messaging/replay`: replay orchestration
- `github.com/PlatformCore/Foundation/messaging/idempotency`: idempotency guard
- `github.com/PlatformCore/Foundation/messaging/example_integration`: adapter/integration example
- `github.com/PlatformCore/Foundation/messaging/audit`: messaging audit events

## 4) Security
- `github.com/PlatformCore/Foundation/security`: security umbrella
- `github.com/PlatformCore/Foundation/security/jwt`: JWT
- `github.com/PlatformCore/Foundation/security/oauth2`: OAuth2
- `github.com/PlatformCore/Foundation/security/hash`: hashing
- `github.com/PlatformCore/Foundation/security/encryption`: encryption
- `github.com/PlatformCore/Foundation/security/permission`: permission checks
- `github.com/PlatformCore/Foundation/security/policy`: policy checks
- `github.com/PlatformCore/Foundation/security/threatdefense`: threat defense rules

## 5) Persistence
- `github.com/driftappdev/persistence`: persistence umbrella
- `github.com/PlatformCore/Foundation/persistence/tx`: transaction helpers
- `github.com/PlatformCore/Foundation/persistence/uow`: unit-of-work
- `github.com/PlatformCore/Foundation/persistence/distlock`: distributed lock

## 6) Observability / Telemetry
- `github.com/driftappdev/observability`: observability umbrella
- `github.com/PlatformCore/Foundation/observability/correlation`: correlation helpers
- `github.com/PlatformCore/Foundation/observability/healthcheck`: liveness/readiness/health
- `github.com/PlatformCore/Foundation/observability/span`: span helpers
- `github.com/PlatformCore/Foundation/observability/tracing`: tracing helpers
- `github.com/PlatformCore/Foundation/observability/performance`: performance metrics helpers
- `github.com/PlatformCore/Foundation/observability/profiler/sentinel`: profiler integration
- `github.com/PlatformCore/Foundation/observability/audit`: observability-side audit
- `github.com/driftappdev/observability/telemetry`: telemetry umbrella (OpenTelemetry-centric)

## 7) Platform / Plugins / Eventing
- `github.com/driftappdev/platform`: platform umbrella
- `github.com/PlatformCore/Foundation/platform/servicemesh`: service mesh integration
- `github.com/PlatformCore/Foundation/plugins`: plugin umbrella
- `github.com/PlatformCore/Foundation/plugins/common`: shared plugin contracts/util
- `github.com/PlatformCore/Foundation/plugins/engine`: plugin engine/runtime
- `github.com/PlatformCore/Foundation/platform/eventbus`: event bus abstractions

## 8) Resilience / Rate limiting
- `github.com/PlatformCore/Foundation/resilience`: resilience umbrella
- `github.com/PlatformCore/Foundation/resilience/retry`: retry
- `github.com/PlatformCore/Foundation/resilience/sanitizer`: sanitize/cleanup
- `github.com/PlatformCore/Foundation/resilience/validate`: validation utilities
- `github.com/PlatformCore/Foundation/resilience/validator`: validator helpers
- `github.com/PlatformCore/Foundation/resilience/pagination`: pagination helpers
- `github.com/PlatformCore/Foundation/resilience/cache`: cache resilience helpers
- `github.com/PlatformCore/Foundation/resilience/circuit`: circuit breaker
- `github.com/PlatformCore/Foundation/ratelimit`: rate limit umbrella
- `github.com/PlatformCore/Foundation/ratelimit/memory_store`: memory store for rate limit
- `github.com/PlatformCore/Foundation/ratelimit/enterprise`: enterprise rate-limit features

## 9) Middleware
- `github.com/driftappdev/middleware`: middleware umbrella
- `github.com/driftappdev/middleware/requestid`: request id middleware
- `github.com/driftappdev/middleware/timeout`: timeout middleware
- `github.com/driftappdev/middleware/logmid/logging-middleware`: logging middleware variant
- `github.com/driftappdev/middleware/adminshield/admin-middleware`: admin shield middleware variant
- `github.com/PlatformCore/Foundation/observability/logging-middleware`: standalone logging middleware
- `github.com/driftappdev/adminshield`: standalone adminshield (namespace à¸žà¸´à¸¡à¸žà¹Œà¸•à¹ˆà¸²à¸‡à¸ˆà¸²à¸ driftappdev)
- `github.com/driftappdev/auth`: standalone auth (namespace à¸žà¸´à¸¡à¸žà¹Œà¸•à¹ˆà¸²à¸‡à¸ˆà¸²à¸ driftappdev)

## 10) Legacy / Compatibility Facades (`go*`)
- `github.com/PlatformCore/Foundation/security/auth_middleware/goauth`
- `github.com/PlatformCore/Foundation/resilience/circuitbreaker/gocircuit`
- `github.com/PlatformCore/Foundation/core/errors/goerror_compat`
- `github.com/PlatformCore/Foundation/observability/logging/gologger`
- `github.com/PlatformCore/Foundation/observability/gometrics`
- `github.com/PlatformCore/Foundation/ratelimit/compat/goratelimit`
- `github.com/PlatformCore/Foundation/resilience/retry/goretry`
- `github.com/PlatformCore/Foundation/security/sanitizer/gosanitizer`
- `github.com/PlatformCore/Foundation/resilience/timeout/gotimeout`
- `github.com/PlatformCore/Foundation/observability/tracing/gotracing`

## 11) Folder Alias Map (à¸Šà¸·à¹ˆà¸­à¹‚à¸Ÿà¸¥à¹€à¸”à¸­à¸£à¹Œ != module path)
à¸à¸¥à¸¸à¹ˆà¸¡à¸™à¸µà¹‰à¹„à¸¡à¹ˆà¹ƒà¸Šà¹ˆà¹‚à¸¡à¸”à¸¹à¸¥à¹€à¸žà¸´à¹ˆà¸¡ à¹à¸•à¹ˆà¹€à¸›à¹‡à¸™à¹‚à¸Ÿà¸¥à¹€à¸”à¸­à¸£à¹Œà¸—à¸µà¹ˆà¸Šà¸µà¹‰à¹„à¸›à¸¢à¸±à¸‡ module path à¹€à¸Šà¸´à¸‡à¹‚à¸”à¹€à¸¡à¸™:
- `audit` -> `github.com/PlatformCore/Foundation/observability/audit`
- `auditv.1` -> `github.com/PlatformCore/Foundation/messaging/audit`
- `cache` -> `github.com/PlatformCore/Foundation/resilience/cache`
- `circuit` -> `github.com/PlatformCore/Foundation/resilience/circuit`
- `retry` -> `github.com/PlatformCore/Foundation/resilience/retry`
- `sanitizer` -> `github.com/PlatformCore/Foundation/resilience/sanitizer`
- `pagination` -> `github.com/PlatformCore/Foundation/resilience/pagination`
- `dlq` -> `github.com/PlatformCore/Foundation/messaging/dlq`
- `inbox` -> `github.com/PlatformCore/Foundation/messaging/inbox`
- `outbox` -> `github.com/PlatformCore/Foundation/messaging/outbox`
- `redrive` -> `github.com/PlatformCore/Foundation/messaging/redrive`
- `replay` -> `github.com/PlatformCore/Foundation/messaging/replay`
- `jwt` -> `github.com/PlatformCore/Foundation/security/jwt`
- `oauth2` -> `github.com/PlatformCore/Foundation/security/oauth2`
- `hash` -> `github.com/PlatformCore/Foundation/security/hash`
- `encryption` -> `github.com/PlatformCore/Foundation/security/encryption`
- `permission` -> `github.com/PlatformCore/Foundation/security/permission`
- `policy` -> `github.com/PlatformCore/Foundation/security/policy`
- `threatdefense` -> `github.com/PlatformCore/Foundation/security/threatdefense`
- `distlock` -> `github.com/PlatformCore/Foundation/persistence/distlock`
- `tx` -> `github.com/PlatformCore/Foundation/persistence/tx`
- `uow` -> `github.com/PlatformCore/Foundation/persistence/uow`
- `servicemesh` -> `github.com/PlatformCore/Foundation/platform/servicemesh`
- `common` -> `github.com/PlatformCore/Foundation/plugins/common`
- `engine` -> `github.com/PlatformCore/Foundation/plugins/engine`
- `enterprise` -> `github.com/PlatformCore/Foundation/ratelimit/enterprise`
- `performance` -> `github.com/PlatformCore/Foundation/observability/performance`
- `sentinelprofiler` -> `github.com/PlatformCore/Foundation/observability/profiler/sentinel`
- `ratelimitX` -> `github.com/PlatformCore/Foundation/ratelimit`

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




