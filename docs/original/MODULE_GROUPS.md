# Module Structure (Mainstream Style)

This repository follows a domain-first Go module layout:
- Main module = domain root (can be installed alone)
- Submodule = optional focused package under that domain
- Submodules are not required unless imported by your app or required by the main module

## Main Modules (With Submodules)

- `github.com/PlatformCore/Foundation/clients`
- `github.com/PlatformCore/Foundation/core/contracts`
- `github.com/PlatformCore/Foundation/core`
- `github.com/PlatformCore/Foundation/orchestration/di`
- `github.com/PlatformCore/Foundation/platform/eventbus`
- `github.com/PlatformCore/Foundation/platform/featureflag`
- `github.com/driftappdev/infra`
- `github.com/driftappdev/observability`
- `github.com/driftappdev/platform`
- `github.com/PlatformCore/Foundation/plugins`
- `github.com/PlatformCore/Foundation/resilience`
- `github.com/PlatformCore/Foundation/runtime`
- `github.com/driftappdev/observability/telemetry`
- `github.com/driftappdev/testing`
- `github.com/driftappdev/foundation/validator`

## Standalone Main Modules (No Submodules)

- `github.com/PlatformCore/Foundation/config`
- `github.com/driftappdev/docs`
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
- `github.com/driftappdev/logmid/logging-middleware`
- `github.com/PlatformCore/Foundation/ratelimit`
- `github.com/PlatformCore/Foundation/core/result_legacy`
- `github.com/PlatformCore/Foundation/security/encryption`
- `github.com/PlatformCore/Foundation/security/hash`
- `github.com/PlatformCore/Foundation/security/jwt`
- `github.com/PlatformCore/Foundation/security/oauth2`
- `github.com/PlatformCore/Foundation/security/secrets`

## Submodules

- `github.com/PlatformCore/Foundation/clients/grpc`
- `github.com/PlatformCore/Foundation/clients/http`
- `github.com/PlatformCore/Foundation/clients/nats`
- `github.com/PlatformCore/Foundation/core/contracts/errors`
- `github.com/PlatformCore/Foundation/core/contracts/pagination`
- `github.com/PlatformCore/Foundation/core/contracts/response`
- `github.com/PlatformCore/Foundation/core/contracts/versioning`
- `github.com/PlatformCore/Foundation/core/constants`
- `github.com/PlatformCore/Foundation/core/context`
- `github.com/PlatformCore/Foundation/core/errors`
- `github.com/PlatformCore/Foundation/core/logger`
- `github.com/PlatformCore/Foundation/core/result`
- `github.com/PlatformCore/Foundation/core/types`
- `github.com/PlatformCore/Foundation/core/utils`
- `github.com/PlatformCore/Foundation/orchestration/di/container`
- `github.com/PlatformCore/Foundation/orchestration/di/module`
- `github.com/PlatformCore/Foundation/orchestration/di/provider`
- `github.com/PlatformCore/Foundation/orchestration/di/registry`
- `github.com/PlatformCore/Foundation/orchestration/di/scope`
- `github.com/PlatformCore/Foundation/platform/eventbus/deadletter`
- `github.com/PlatformCore/Foundation/platform/eventbus/envelope`
- `github.com/PlatformCore/Foundation/platform/eventbus/headers`
- `github.com/PlatformCore/Foundation/platform/eventbus/idempotency`
- `github.com/PlatformCore/Foundation/platform/eventbus/publisher`
- `github.com/PlatformCore/Foundation/platform/eventbus/registry`
- `github.com/PlatformCore/Foundation/platform/eventbus/retry`
- `github.com/PlatformCore/Foundation/platform/eventbus/serializer`
- `github.com/PlatformCore/Foundation/platform/eventbus/subscriber`
- `github.com/PlatformCore/Foundation/platform/featureflag/cache`
- `github.com/PlatformCore/Foundation/platform/featureflag/client`
- `github.com/PlatformCore/Foundation/platform/featureflag/evaluator`
- `github.com/PlatformCore/Foundation/platform/featureflag/provider`
- `github.com/PlatformCore/Foundation/platform/featureflag/types`
- `github.com/driftappdev/infra/backoff`
- `github.com/driftappdev/infra/bulkhead`
- `github.com/driftappdev/infra/cache`
- `github.com/driftappdev/infra/circuit`
- `github.com/driftappdev/infra/clock`
- `github.com/driftappdev/infra/retry`
- `github.com/PlatformCore/Foundation/observability/correlation`
- `github.com/PlatformCore/Foundation/observability/healthcheck`
- `github.com/PlatformCore/Foundation/observability/logging`
- `github.com/PlatformCore/Foundation/observability/profiler`
- `github.com/PlatformCore/Foundation/observability/span`
- `github.com/driftappdev/observability/trace`
- `github.com/PlatformCore/Foundation/observability/tracing`
- `github.com/PlatformCore/Foundation/clients/platform`
- `github.com/PlatformCore/Foundation/platform/container`
- `github.com/PlatformCore/Foundation/platform/evaluator`
- `github.com/driftappdev/platform/hooks`
- `github.com/driftappdev/platform/loader`
- `github.com/driftappdev/platform/provider`
- `github.com/driftappdev/platform/registry`
- `github.com/PlatformCore/Foundation/platform/versioning`
- `github.com/PlatformCore/Foundation/plugins/hooks`
- `github.com/PlatformCore/Foundation/plugins/loader`
- `github.com/PlatformCore/Foundation/plugins/manifest`
- `github.com/PlatformCore/Foundation/plugins/registry`
- `github.com/PlatformCore/Foundation/resilience/cache`
- `github.com/PlatformCore/Foundation/resilience/circuit`
- `github.com/PlatformCore/Foundation/resilience/pagination`
- `github.com/PlatformCore/Foundation/resilience/retry`
- `github.com/PlatformCore/Foundation/resilience/sanitizer`
- `github.com/PlatformCore/Foundation/resilience/schema`
- `github.com/PlatformCore/Foundation/resilience/serializer`
- `github.com/PlatformCore/Foundation/resilience/validate`
- `github.com/PlatformCore/Foundation/resilience/validator`
- `github.com/PlatformCore/Foundation/runtime/health`
- `github.com/PlatformCore/Foundation/runtime/lifecycle`
- `github.com/PlatformCore/Foundation/runtime/shutdown`
- `github.com/PlatformCore/Foundation/runtime/signals`
- `github.com/PlatformCore/Foundation/observability/telemetry/baggage`
- `github.com/PlatformCore/Foundation/observability/telemetry/correlation`
- `github.com/PlatformCore/Foundation/observability/telemetry/trace`
- `github.com/driftappdev/testing/fixtures`
- `github.com/driftappdev/testing/mocks`
- `github.com/driftappdev/testing/testutil`
- `github.com/PlatformCore/Foundation/validation/binding`
- `github.com/PlatformCore/Foundation/validation/schema`

## Install Examples

```bash
go get github.com/PlatformCore/Foundation/clients@latest
go get github.com/PlatformCore/Foundation/core/contracts@latest
go get github.com/PlatformCore/Foundation/core@latest
go get github.com/PlatformCore/Foundation/clients/grpc@latest
go get github.com/PlatformCore/Foundation/clients/http@latest
go get github.com/PlatformCore/Foundation/clients/nats@latest
```



