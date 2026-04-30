# Libpackage Enterprise Package Audit

Generated: 2026-04-28 04:32:23 +07:00

Total modules (go.mod): 83

## Duplicate/Overlap Clusters (heuristic)
- middleware stack (18 modules)
  - github.com/PlatformCore/libpackage/observability/logging-middleware, github.com/driftappdev/middleware, github.com/driftappdev/middleware/adminshield/admin-middleware, github.com/driftappdev/middleware/requestid, github.com/driftappdev/middleware/timeout, github.com/enterprise/middleware, github.com/PlatformCore/libpackage/middleware/clock, github.com/enterprise/middleware/cmd/example, github.com/PlatformCore/libpackage/middleware/event, github.com/enterprise/middleware/grpc, github.com/enterprise/middleware/http, github.com/PlatformCore/libpackage/middleware/ids, github.com/enterprise/middleware/mnt/user-data/outputs/obslib/pkg/middleware/event, github.com/PlatformCore/libpackage/middleware/pool, github.com/PlatformCore/libpackage/middleware/propagation, github.com/PlatformCore/libpackage/middleware/registry, github.com/enterprise/middleware/telemetry, github.com/PlatformCore/libpackage/middleware/trace
- tracing/telemetry stack (12 modules)
  - github.com/PlatformCore/libpackage/observability/tracing/gotracing, github.com/driftappdev/observability, github.com/PlatformCore/libpackage/observability/audit, github.com/PlatformCore/libpackage/observability/correlation, github.com/PlatformCore/libpackage/observability/healthcheck, github.com/PlatformCore/libpackage/observability/performance, github.com/PlatformCore/libpackage/observability/profiler/sentinel, github.com/PlatformCore/libpackage/observability/span, github.com/PlatformCore/libpackage/observability/tracing, github.com/driftappdev/observability/telemetry, github.com/enterprise/middleware/telemetry, github.com/PlatformCore/libpackage/middleware/trace
- security/auth stack (10 modules)
  - github.com/driftappdev/auth, github.com/PlatformCore/libpackage/security/auth_middleware/goauth, github.com/PlatformCore/libpackage/security, github.com/PlatformCore/libpackage/security/encryption, github.com/PlatformCore/libpackage/security/hash, github.com/PlatformCore/libpackage/security/jwt, github.com/PlatformCore/libpackage/security/oauth2, github.com/PlatformCore/libpackage/security/permission, github.com/PlatformCore/libpackage/security/policy, github.com/PlatformCore/libpackage/security/threatdefense
- rate-limit stack (4 modules)
  - github.com/PlatformCore/libpackage/ratelimit/compat/goratelimit, github.com/PlatformCore/libpackage/ratelimit, github.com/PlatformCore/libpackage/ratelimit/enterprise, github.com/PlatformCore/libpackage/ratelimit/memory_store
- resilience stack (4 modules)
  - github.com/PlatformCore/libpackage/resilience/circuitbreaker/gocircuit, github.com/PlatformCore/libpackage/tools/lib_word/flowguard-ultimate, github.com/PlatformCore/libpackage/resilience, github.com/PlatformCore/libpackage/resilience/validator
- timeout/deadline stack (2 modules)
  - github.com/PlatformCore/libpackage/resilience/timeout/gotimeout, github.com/driftappdev/middleware/timeout
- validation stack (2 modules)
  - github.com/PlatformCore/libpackage/resilience/validator, github.com/driftappdev/foundation/validator
- logging stack (2 modules)
  - github.com/PlatformCore/libpackage/observability/logging/gologger, github.com/PlatformCore/libpackage/observability/logging-middleware

## Module Inventory
### github.com/driftappdev/adminshield
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\adminshield
- Role: runtime/platform/integration
- Go files: 1
- Files:
  - adminshield.go

### github.com/driftappdev/auth
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\auth
- Role: auth/security
- Go files: 4
- Files:
  - auth.go
  - auth-middleware\auth.middleware.go
  - gin\gin.go
  - grpc\grpc.go

### github.com/driftappdev
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage
- Role: general utility
- Go files: 2
- Files:
  - schema\schema.go
  - workflow\workflow.go

### github.com/PlatformCore/libpackage/clients
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\client
- Role: general utility
- Go files: 8
- Files:
  - grpc\client.go
  - grpc\interceptor.go
  - http\client.go
  - http\options.go
  - http\retry.go
  - http\timeout.go
  - nats\publisher.go
  - nats\subscriber.go

### github.com/PlatformCore/libpackage/config
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\config
- Role: general utility
- Go files: 5
- Files:
  - defaults\defaults.go
  - env\env.go
  - file\file.go
  - loader\loader.go
  - watcher\watcher.go

### github.com/PlatformCore/libpackage/core/contracts
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\contracts
- Role: general utility
- Go files: 4
- Files:
  - errors\public.go
  - pagination\pagination.go
  - response\envelope.go
  - versioning\version.go

### github.com/PlatformCore/libpackage/core
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\core
- Role: general utility
- Go files: 21
- Files:
  - constants\constants.go
  - context\context.go
  - context\keys.go
  - context\values.go
  - errors\advanced_base.go
  - errors\base.go
  - errors\code.go
  - errors\errors.go
  - errors\wrap.go
  - logger\advanced.go
  - logger\field.go
  - logger\logger.go
  - logger\options.go
  - result\result.go
  - types\nullable.go
  - types\types.go
  - utils\pointer.go
  - utils\slice.go
  - utils\string.go
  - utils\time.go
  - utils\utils.go

### github.com/PlatformCore/libpackage/orchestration/di
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\di
- Role: runtime/platform/integration
- Go files: 5
- Files:
  - container\container.go
  - module\module.go
  - provider\provider.go
  - registry\registry.go
  - scope\scope.go

### github.com/PlatformCore/libpackage/platform/eventbus
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\eventbus
- Role: messaging/event pipeline
- Go files: 10
- Files:
  - codec\codec.go
  - deadletter\dlq.go
  - envelope\envelope.go
  - headers\headers.go
  - idempotency\store.go
  - publisher\publisher.go
  - registry\registry.go
  - retry\retry.go
  - serializer\serializer.go
  - subscriber\subscriber.go

### github.com/PlatformCore/libpackage/platform/featureflag
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\featureflag
- Role: general utility
- Go files: 13
- Files:
  - cache\cache.go
  - client\client.go
  - client\memory.go
  - enterprise\featureflag.go
  - evaluator\advanced.go
  - evaluator\evaluator.go
  - provider\file.go
  - provider\memory.go
  - provider\redis.go
  - types\flag.go
  - types\rule.go
  - types\target.go
  - types\variant.go

### github.com/PlatformCore/libpackage/security/auth_middleware/goauth
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\goauth
- Role: auth/security
- Go files: 1
- Files:
  - auth.middleware.go

### github.com/PlatformCore/libpackage/resilience/circuitbreaker/gocircuit
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\gocircuit
- Role: resilience primitives
- Go files: 1
- Files:
  - circuitBreaker.go

### github.com/PlatformCore/libpackage/core/errors/goerror_compat
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\goerror
- Role: general utility
- Go files: 1
- Files:
  - error.go

### github.com/PlatformCore/libpackage/observability/logging/gologger
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\gologger
- Role: observability / telemetry
- Go files: 1
- Files:
  - logger.go

### github.com/PlatformCore/libpackage/observability/gometrics
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\gometrics
- Role: observability / telemetry
- Go files: 1
- Files:
  - metrics.go

### github.com/PlatformCore/libpackage/ratelimit/compat/goratelimit
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\goratelimit
- Role: rate limiting / throttling
- Go files: 6
- Files:
  - advanced_engine.go
  - advanced_memory_store.go
  - advanced_redis_store.go
  - advanced_types.go
  - framework_middleware.go
  - rateLimit.middleware.go

### github.com/PlatformCore/libpackage/resilience/retry/goretry
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\goretry
- Role: retry/backoff
- Go files: 1
- Files:
  - retry.go

### github.com/PlatformCore/libpackage/security/sanitizer/gosanitizer
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\gosanitizer
- Role: auth/security
- Go files: 1
- Files:
  - sanitizer.go

### github.com/PlatformCore/libpackage/resilience/timeout/gotimeout
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\gotimeout
- Role: timeout / deadlines
- Go files: 1
- Files:
  - timeout.go

### github.com/PlatformCore/libpackage/observability/tracing/gotracing
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\gotracing
- Role: observability / telemetry
- Go files: 1
- Files:
  - tracing.go

### github.com/PlatformCore/libpackage/tools/lib_word
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\lib_word
- Role: general utility
- Go files: 13
- Files:
  - adaptive_concurrency.go
  - backpressure.go
  - batcher.go
  - bulkhead.go
  - checkpoint.go
  - deadline.go
  - health_supervisor.go
  - limiter.go
  - load_shedder.go
  - priority_queue.go
  - retry_backoff.go
  - token_bucket.go
  - workqueue.go

### github.com/PlatformCore/libpackage/tools/lib_word/cmd/example
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\lib_word\cmd\example
- Role: general utility
- Go files: 1
- Files:
  - main.go

### github.com/PlatformCore/libpackage/tools/lib_word/flowguard-ultimate
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\lib_word\flowguard-ultimate
- Role: resilience primitives
- Go files: 10
- Files:
  - adaptive_retry.go
  - cache_stampede_protection.go
  - congestion_control.go
  - distributed_rate_limit.go
  - hedged_request.go
  - quorum.go
  - request_coalescing.go
  - shadow_traffic.go
  - state_synchronizer.go
  - traffic_mirroring.go

### github.com/PlatformCore/libpackage/observability/logging-middleware
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\logging-middleware
- Role: observability / telemetry
- Go files: 1
- Files:
  - logging.middleware.go

### github.com/PlatformCore/libpackage/messaging/audit
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\auditv.1
- Role: runtime/platform/integration
- Go files: 1
- Files:
  - audit.go

### github.com/PlatformCore/libpackage/messaging/dlq
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\dlq
- Role: messaging/event pipeline
- Go files: 1
- Files:
  - dlq.go

### github.com/PlatformCore/libpackage/messaging/example_integration
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\example_integration
- Role: general utility
- Go files: 1
- Files:
  - example_inte.go

### github.com/PlatformCore/libpackage/messaging/idempotency
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\idempotency
- Role: general utility
- Go files: 1
- Files:
  - idempotency.go

### github.com/PlatformCore/libpackage/messaging/inbox
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\inbox
- Role: messaging/event pipeline
- Go files: 1
- Files:
  - inbox.go

### github.com/PlatformCore/libpackage/messaging/outbox
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\outbox
- Role: messaging/event pipeline
- Go files: 1
- Files:
  - outbox.go

### github.com/PlatformCore/libpackage/messaging/redrive
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\redrive
- Role: messaging/event pipeline
- Go files: 1
- Files:
  - redrive.go

### github.com/PlatformCore/libpackage/messaging/replay
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\replay
- Role: messaging/event pipeline
- Go files: 6
- Files:
  - checkpoint.go
  - errors.go
  - range.go
  - runner.go
  - selector.go
  - service.go

### github.com/driftappdev/middleware
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware
- Role: http/grpc middleware
- Go files: 16
- Files:
  - cors\cors.go
  - cors\gin.go
  - logging\gin\gin.go
  - logging\grpc\grpc.go
  - ratelimit\extractor\extractor.go
  - ratelimit\gin\gin.go
  - ratelimit\grpc\grpc.go
  - ratelimit\headers\headers.go
  - ratelimit\ratelimit-middleware\rateLimit.middleware.go
  - recovery\gin.go
  - recovery\gin\gin.go
  - recovery\grpc.go
  - recovery\grpc\grpc.go
  - recovery\recovery.go
  - tracing\gin\gin.go
  - tracing\grpc\grpc.go

### github.com/driftappdev/middleware/adminshield/admin-middleware
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\adminshield\admin-middleware
- Role: http/grpc middleware
- Go files: 1
- Files:
  - admin.middleware.go

### github.com/driftappdev/middleware/requestid
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\requestid
- Role: http/grpc middleware
- Go files: 5
- Files:
  - gin.go
  - gin\gin.go
  - grpc.go
  - grpc\grpc.go
  - requestid.go

### github.com/driftappdev/middleware/timeout
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\timeout
- Role: timeout / deadlines
- Go files: 5
- Files:
  - gin.go
  - gin\gin.go
  - grpc.go
  - grpc\grpc.go
  - timeout.go

### github.com/driftappdev/observability
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\observability
- Role: observability / telemetry
- Go files: 8
- Files:
  - logging\context.go
  - logging\fields.go
  - logging\logging.go
  - metrics\metrics.go
  - metrics\registry\registry.go
  - profiler\pprof.go
  - profiler\profiler.go
  - trace\trace.go

### github.com/PlatformCore/libpackage/observability/audit
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\audit
- Role: observability / telemetry
- Go files: 1
- Files:
  - audit.go

### github.com/PlatformCore/libpackage/observability/correlation
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\observability\correlation
- Role: observability / telemetry
- Go files: 1
- Files:
  - correlation.go

### github.com/PlatformCore/libpackage/observability/healthcheck
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\observability\healthcheck
- Role: observability / telemetry
- Go files: 3
- Files:
  - healthcheck.go
  - liveness.go
  - readiness.go

### github.com/PlatformCore/libpackage/observability/performance
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\performance
- Role: observability / telemetry
- Go files: 1
- Files:
  - performance.go

### github.com/PlatformCore/libpackage/observability/profiler/sentinel
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\sentinelprofiler
- Role: observability / telemetry
- Go files: 1
- Files:
  - profiler.go

### github.com/PlatformCore/libpackage/observability/span
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\observability\span
- Role: observability / telemetry
- Go files: 1
- Files:
  - span.go

### github.com/PlatformCore/libpackage/observability/tracing
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\observability\tracing
- Role: observability / telemetry
- Go files: 3
- Files:
  - propagator.go
  - sampler.go
  - tracing.go

### github.com/driftappdev/persistence
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\persistence
- Role: data access / transaction
- Go files: 0

### github.com/PlatformCore/libpackage/persistence/distlock
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\distlock
- Role: data access / transaction
- Go files: 1
- Files:
  - distlock.go

### github.com/PlatformCore/libpackage/persistence/tx
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\tx
- Role: data access / transaction
- Go files: 1
- Files:
  - tx.go

### github.com/PlatformCore/libpackage/persistence/uow
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\uow
- Role: data access / transaction
- Go files: 1
- Files:
  - uow.go

### github.com/driftappdev/platform
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\platform
- Role: runtime/platform/integration
- Go files: 4
- Files:
  - client\client.go
  - container\container.go
  - evaluator\evaluator.go
  - versioning\versioning.go

### github.com/PlatformCore/libpackage/platform/servicemesh
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\servicemesh
- Role: runtime/platform/integration
- Go files: 1
- Files:
  - servicemesh.go

### github.com/PlatformCore/libpackage/plugins
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\plugins
- Role: runtime/platform/integration
- Go files: 5
- Files:
  - hooks\hooks.go
  - loader\loader.go
  - manifest\manifest.go
  - provider\provider.go
  - registry\registry.go

### github.com/PlatformCore/libpackage/plugins/common
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\common
- Role: runtime/platform/integration
- Go files: 1
- Files:
  - common.go

### github.com/PlatformCore/libpackage/plugins/engine
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\engine
- Role: runtime/platform/integration
- Go files: 1
- Files:
  - engine.go

### github.com/PlatformCore/libpackage/ratelimit
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\ratelimitX
- Role: rate limiting / throttling
- Go files: 9
- Files:
  - errors\errors.go
  - key\key.go
  - limiter\limiter.go
  - lua\lua.go
  - options\options.go
  - policy\policy.go
  - ratelimit.go
  - redis_store\redis_store.go
  - result\result.go

### github.com/PlatformCore/libpackage/ratelimit/enterprise
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\enterprise
- Role: rate limiting / throttling
- Go files: 1
- Files:
  - ratelimit.go

### github.com/PlatformCore/libpackage/ratelimit/memory_store
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\ratelimitX\memory_store
- Role: rate limiting / throttling
- Go files: 1
- Files:
  - memory_store.go

### github.com/PlatformCore/libpackage/resilience
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\resilience
- Role: resilience primitives
- Go files: 0

### github.com/PlatformCore/libpackage/resilience/validator
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\resilience\validator
- Role: resilience primitives
- Go files: 1
- Files:
  - validator.go

### github.com/PlatformCore/libpackage/core/result_legacy
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\result
- Role: general utility
- Go files: 1
- Files:
  - result.go

### github.com/PlatformCore/libpackage/runtime
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\runtime
- Role: runtime/platform/integration
- Go files: 4
- Files:
  - health\health.go
  - lifecycle\lifecycle.go
  - shutdown\shutdown.go
  - signals\signals.go

### github.com/PlatformCore/libpackage/security
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\security
- Role: auth/security
- Go files: 1
- Files:
  - security.go

### github.com/PlatformCore/libpackage/security/encryption
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\encryption
- Role: auth/security
- Go files: 1
- Files:
  - encryption.go

### github.com/PlatformCore/libpackage/security/hash
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\hash
- Role: auth/security
- Go files: 1
- Files:
  - hash.go

### github.com/PlatformCore/libpackage/security/jwt
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\jwt
- Role: auth/security
- Go files: 1
- Files:
  - jwt.go

### github.com/PlatformCore/libpackage/security/oauth2
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\oauth2
- Role: auth/security
- Go files: 2
- Files:
  - client.go
  - oauth2.go

### github.com/PlatformCore/libpackage/security/permission
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\permission
- Role: auth/security
- Go files: 1
- Files:
  - permission.go

### github.com/PlatformCore/libpackage/security/policy
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\policy
- Role: auth/security
- Go files: 1
- Files:
  - policy.go

### github.com/PlatformCore/libpackage/security/threatdefense
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\threatdefense
- Role: auth/security
- Go files: 1
- Files:
  - threat_defense.go

### github.com/driftappdev/observability/telemetry
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\telemetry
- Role: observability / telemetry
- Go files: 3
- Files:
  - baggage\baggage.go
  - correlation\context.go
  - trace\provider.go

### github.com/driftappdev/foundation/validator
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\validator
- Role: validation/schema
- Go files: 5
- Files:
  - binding\grpc.go
  - binding\http.go
  - schema\engine.go
  - schema\errors.go
  - schema\rules.go

### github.com/enterprise/middleware
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware
- Role: http/grpc middleware
- Go files: 0

### github.com/PlatformCore/libpackage/middleware/clock
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\clock
- Role: http/grpc middleware
- Go files: 1
- Files:
  - clock.go

### github.com/enterprise/middleware/cmd/example
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\cmd\example
- Role: http/grpc middleware
- Go files: 1
- Files:
  - main.go

### github.com/PlatformCore/libpackage/middleware/event
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\event
- Role: http/grpc middleware
- Go files: 1
- Files:
  - middleware.go

### github.com/enterprise/middleware/grpc
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\grpc
- Role: http/grpc middleware
- Go files: 1
- Files:
  - interceptor.go

### github.com/enterprise/middleware/http
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\http
- Role: http/grpc middleware
- Go files: 17
- Files:
  - auth.go
  - cache.go
  - circuit_breaker.go
  - compress.go
  - cors.go
  - example_test.go
  - health.go
  - logger.go
  - metrics.go
  - middleware.go
  - ratelimit.go
  - recovery.go
  - requestid.go
  - security.go
  - timeout.go
  - tracing.go
  - validator.go

### github.com/PlatformCore/libpackage/middleware/ids
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\ids
- Role: http/grpc middleware
- Go files: 1
- Files:
  - ids.go

### github.com/enterprise/middleware/mnt/user-data/outputs/obslib/pkg/middleware/event
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\mnt\user-data\outputs\obslib\pkg\middleware\event
- Role: http/grpc middleware
- Go files: 1
- Files:
  - middleware.go

### github.com/PlatformCore/libpackage/middleware/pool
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\pool
- Role: http/grpc middleware
- Go files: 1
- Files:
  - bytespool.go

### github.com/PlatformCore/libpackage/middleware/propagation
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\propagation
- Role: http/grpc middleware
- Go files: 1
- Files:
  - propagation.go

### github.com/PlatformCore/libpackage/middleware/registry
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\registry
- Role: http/grpc middleware
- Go files: 1
- Files:
  - registry.go

### github.com/enterprise/middleware/telemetry
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\telemetry
- Role: observability / telemetry
- Go files: 1
- Files:
  - provider.go

### github.com/PlatformCore/libpackage/middleware/trace
- Path: C:\Users\AdminWC\Dift App project\Project-Production-Ready\libpackage\middleware\push_code_middleware\trace
- Role: observability / telemetry
- Go files: 1
- Files:
  - span.go





