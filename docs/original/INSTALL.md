# libpackage Multi-Module Installation Guide

Repository: `github.com/driftappdev`

This repository is configured as **separate Go modules** (not one combined module).
Install only the module you need.

## Example install commands

```bash
go get github.com/PlatformCore/libpackage/core/types@latest
go get github.com/PlatformCore/libpackage/core/result_legacy@latest
go get github.com/PlatformCore/libpackage/core/context@latest
go get github.com/PlatformCore/libpackage/core/constants@latest
go get github.com/PlatformCore/libpackage/core/errors@latest
go get github.com/PlatformCore/libpackage/core/logger@latest
go get github.com/driftappdev/logmid/logging-middleware@latest
go get github.com/PlatformCore/libpackage/security/jwt@latest
```


## Other available modules

- `github.com/PlatformCore/libpackage/security/jwt`
- `github.com/PlatformCore/libpackage/security/oauth2`
- `github.com/PlatformCore/libpackage/security/hash`
- `github.com/PlatformCore/libpackage/runtime/lifecycle`
- `github.com/PlatformCore/libpackage/runtime/shutdown`
- `github.com/PlatformCore/libpackage/runtime/health`
- `github.com/PlatformCore/libpackage/security/auth_middleware/goauth`
- `github.com/PlatformCore/libpackage/resilience/cache`
- `github.com/PlatformCore/libpackage/resilience/circuitbreaker/gocircuit`
- `github.com/PlatformCore/libpackage/core/errors/goerror_compat`
- `github.com/PlatformCore/libpackage/observability/logging/gologger`
- `github.com/PlatformCore/libpackage/observability/gometrics`
- `github.com/PlatformCore/libpackage/resilience/pagination`
- `github.com/PlatformCore/libpackage/ratelimit/compat/goratelimit`
- `github.com/PlatformCore/libpackage/resilience/retry/goretry`
- `github.com/PlatformCore/libpackage/security/sanitizer/gosanitizer`
- `github.com/PlatformCore/libpackage/resilience/timeout/gotimeout`
- `github.com/PlatformCore/libpackage/observability/tracing/gotracing`
- `github.com/PlatformCore/libpackage/resilience/validate`
- `github.com/PlatformCore/libpackage/resilience/validator`
- `github.com/PlatformCore/libpackage/security/encryption`
### Persistence
- `github.com/PlatformCore/libpackage/persistence/tx`
- `github.com/PlatformCore/libpackage/persistence/uow`

### Messaging
- `github.com/PlatformCore/libpackage/messaging/outbox`
- `github.com/PlatformCore/libpackage/messaging/inbox`
- `github.com/PlatformCore/libpackage/messaging/dlq`
- `github.com/PlatformCore/libpackage/messaging/redrive`
- `github.com/PlatformCore/libpackage/messaging/replay`
## Local development in this repo

- This repo now uses `go.work` to develop multiple modules together.
- Each module has its own `go.mod`.

## Tagging when publishing

For submodules, use tags by module path prefix, for example:

```bash
git tag types/v1.0.0
git tag result/v1.0.0
git tag security/jwt/v1.0.0
git tag logmid/logging-middleware/v1.0.0
```

Then push tags:

```bash
git push origin --tags
```



