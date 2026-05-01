# Module Structure (Mainstream Style)

This repository follows a domain-first Go module layout:
- Main module = domain root (can be installed alone)
- Submodule = optional focused package under that domain
- Submodules are not required unless imported by your app or required by the main module

## Main Modules (With Submodules)

- `github.com/PlatformCore/libpackage/clients`
- `github.com/PlatformCore/libpackage/core/contracts`
- `github.com/PlatformCore/libpackage/core`
- `github.com/PlatformCore/libpackage/orchestration/di`
- `github.com/PlatformCore/libpackage/platform/eventbus`
- `github.com/PlatformCore/libpackage/platform/featureflag`
- `github.com/driftappdev/infra`
- `github.com/driftappdev/observability`
- `github.com/driftappdev/platform`
- `github.com/PlatformCore/libpackage/plugins`
- `github.com/PlatformCore/libpackage/resilience`
- `github.com/PlatformCore/libpackage/runtime`
- `github.com/driftappdev/observability/telemetry`
- `github.com/driftappdev/testing`
- `github.com/driftappdev/foundation/validator`

## Standalone Main Modules (No Submodules)

- `github.com/PlatformCore/libpackage/config`
- `github.com/driftappdev/docs`
- `github.com/PlatformCore/libpackage/security/auth_middleware/goauth`
- `github.com/PlatformCore/libpackage/resilience/circuitbreaker/gocircuit`
- `github.com/PlatformCore/libpackage/core/errors/goerror_compat`
- `github.com/PlatformCore/libpackage/observability/logging/gologger`
- `github.com/PlatformCore/libpackage/observability/gometrics`
- `github.com/PlatformCore/libpackage/resilience/ratelimit/compat/goratelimit`
- `github.com/PlatformCore/libpackage/resilience/retry/goretry`
- `github.com/PlatformCore/libpackage/security/sanitizer/gosanitizer`
- `github.com/PlatformCore/libpackage/resilience/timeout/gotimeout`
- `github.com/PlatformCore/libpackage/observability/tracing/gotracing`
- `github.com/driftappdev/logmid/logging-middleware`
- `github.com/PlatformCore/libpackage/resilience/ratelimit`
- `github.com/PlatformCore/libpackage/core/result_legacy`
- `github.com/PlatformCore/libpackage/security/encryption`
- `github.com/PlatformCore/libpackage/security/hash`
- `github.com/PlatformCore/libpackage/security/jwt`
- `github.com/PlatformCore/libpackage/security/oauth2`
- `github.com/PlatformCore/libpackage/security/secrets`

## Submodules

- `github.com/PlatformCore/libpackage/clients/grpc`
- `github.com/PlatformCore/libpackage/clients/http`
- `github.com/PlatformCore/libpackage/clients/nats`
- `github.com/PlatformCore/libpackage/core/contracts/errors`
- `github.com/PlatformCore/libpackage/core/contracts/pagination`
- `github.com/PlatformCore/libpackage/core/contracts/response`
- `github.com/PlatformCore/libpackage/core/contracts/versioning`
- `github.com/PlatformCore/libpackage/core/constants`
- `github.com/PlatformCore/libpackage/core/context`
- `github.com/PlatformCore/libpackage/core/errors`
- `github.com/PlatformCore/libpackage/observability/logging`
- `github.com/PlatformCore/libpackage/core/result`
- `github.com/PlatformCore/libpackage/core/types`
- `github.com/PlatformCore/libpackage/core/utils`
- `github.com/PlatformCore/libpackage/orchestration/di/container`
- `github.com/PlatformCore/libpackage/orchestration/di/module`
- `github.com/PlatformCore/libpackage/orchestration/di/provider`
- `github.com/PlatformCore/libpackage/orchestration/di/registry`
- `github.com/PlatformCore/libpackage/orchestration/di/scope`
- `github.com/PlatformCore/libpackage/platform/eventbus/deadletter`
- `github.com/PlatformCore/libpackage/platform/eventbus/envelope`
- `github.com/PlatformCore/libpackage/platform/eventbus/headers`
- `github.com/PlatformCore/libpackage/platform/eventbus/idempotency`
- `github.com/PlatformCore/libpackage/platform/eventbus/publisher`
- `github.com/PlatformCore/libpackage/platform/eventbus/registry`
- `github.com/PlatformCore/libpackage/platform/eventbus/retry`
- `github.com/PlatformCore/libpackage/platform/eventbus/serializer`
- `github.com/PlatformCore/libpackage/platform/eventbus/subscriber`
- `github.com/PlatformCore/libpackage/platform/featureflag/cache`
- `github.com/PlatformCore/libpackage/platform/featureflag/client`
- `github.com/PlatformCore/libpackage/platform/featureflag/evaluator`
- `github.com/PlatformCore/libpackage/platform/featureflag/provider`
- `github.com/PlatformCore/libpackage/platform/featureflag/types`
- `github.com/driftappdev/infra/backoff`
- `github.com/driftappdev/infra/bulkhead`
- `github.com/driftappdev/infra/cache`
- `github.com/driftappdev/infra/circuit`
- `github.com/driftappdev/infra/clock`
- `github.com/driftappdev/infra/retry`
- `github.com/PlatformCore/libpackage/observability/correlation`
- `github.com/PlatformCore/libpackage/observability/healthcheck`
- `github.com/PlatformCore/libpackage/observability/logging`
- `github.com/PlatformCore/libpackage/observability/profiler`
- `github.com/PlatformCore/libpackage/observability/span`
- `github.com/driftappdev/observability/trace`
- `github.com/PlatformCore/libpackage/observability/tracing`
- `github.com/PlatformCore/libpackage/clients/platform`
- `github.com/PlatformCore/libpackage/platform/container`
- `github.com/PlatformCore/libpackage/platform/evaluator`
- `github.com/driftappdev/platform/hooks`
- `github.com/driftappdev/platform/loader`
- `github.com/driftappdev/platform/provider`
- `github.com/driftappdev/platform/registry`
- `github.com/PlatformCore/libpackage/platform/versioning`
- `github.com/PlatformCore/libpackage/plugins/hooks`
- `github.com/PlatformCore/libpackage/plugins/loader`
- `github.com/PlatformCore/libpackage/plugins/manifest`
- `github.com/PlatformCore/libpackage/plugins/registry`
- `github.com/PlatformCore/libpackage/resilience/cache`
- `github.com/PlatformCore/libpackage/resilience/circuit`
- `github.com/PlatformCore/libpackage/resilience/pagination`
- `github.com/PlatformCore/libpackage/resilience/retry`
- `github.com/PlatformCore/libpackage/resilience/sanitizer`
- `github.com/PlatformCore/libpackage/resilience/schema`
- `github.com/PlatformCore/libpackage/resilience/serializer`
- `github.com/PlatformCore/libpackage/resilience/validate`
- `github.com/PlatformCore/libpackage/resilience/validator`
- `github.com/PlatformCore/libpackage/runtime/health`
- `github.com/PlatformCore/libpackage/runtime/lifecycle`
- `github.com/PlatformCore/libpackage/runtime/shutdown`
- `github.com/PlatformCore/libpackage/runtime/signals`
- `github.com/PlatformCore/libpackage/observability/telemetry/baggage`
- `github.com/PlatformCore/libpackage/observability/telemetry/correlation`
- `github.com/PlatformCore/libpackage/observability/telemetry/trace`
- `github.com/driftappdev/testing/fixtures`
- `github.com/driftappdev/testing/mocks`
- `github.com/driftappdev/testing/testutil`
- `github.com/PlatformCore/libpackage/validation/binding`
- `github.com/PlatformCore/libpackage/validation/schema`

## Install Examples

```bash
go get github.com/PlatformCore/libpackage/clients@latest
go get github.com/PlatformCore/libpackage/core/contracts@latest
go get github.com/PlatformCore/libpackage/core@latest
go get github.com/PlatformCore/libpackage/clients/grpc@latest
go get github.com/PlatformCore/libpackage/clients/http@latest
go get github.com/PlatformCore/libpackage/clients/nats@latest
```



