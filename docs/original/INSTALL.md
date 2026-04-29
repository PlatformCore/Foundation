# libpackage Multi-Module Installation Guide

Repository: `github.com/driftappdev`

This repository is configured as **separate Go modules** (not one combined module).
Install only the module you need.

## Example install commands

```bash
go get github.com/PlatformCore/Foundation/core/types@latest
go get github.com/PlatformCore/Foundation/core/result_legacy@latest
go get github.com/PlatformCore/Foundation/core/context@latest
go get github.com/PlatformCore/Foundation/core/constants@latest
go get github.com/PlatformCore/Foundation/core/errors@latest
go get github.com/PlatformCore/Foundation/core/logger@latest
go get github.com/driftappdev/logmid/logging-middleware@latest
go get github.com/PlatformCore/Foundation/security/jwt@latest
```


## Other available modules

- `github.com/PlatformCore/Foundation/security/jwt`
- `github.com/PlatformCore/Foundation/security/oauth2`
- `github.com/PlatformCore/Foundation/security/hash`
- `github.com/PlatformCore/Foundation/runtime/lifecycle`
- `github.com/PlatformCore/Foundation/runtime/shutdown`
- `github.com/PlatformCore/Foundation/runtime/health`
- `github.com/PlatformCore/Foundation/security/auth_middleware/goauth`
- `github.com/PlatformCore/Foundation/resilience/cache`
- `github.com/PlatformCore/Foundation/resilience/circuitbreaker/gocircuit`
- `github.com/PlatformCore/Foundation/core/errors/goerror_compat`
- `github.com/PlatformCore/Foundation/observability/logging/gologger`
- `github.com/PlatformCore/Foundation/observability/gometrics`
- `github.com/PlatformCore/Foundation/resilience/pagination`
- `github.com/PlatformCore/Foundation/ratelimit/compat/goratelimit`
- `github.com/PlatformCore/Foundation/resilience/retry/goretry`
- `github.com/PlatformCore/Foundation/security/sanitizer/gosanitizer`
- `github.com/PlatformCore/Foundation/resilience/timeout/gotimeout`
- `github.com/PlatformCore/Foundation/observability/tracing/gotracing`
- `github.com/PlatformCore/Foundation/resilience/validate`
- `github.com/PlatformCore/Foundation/resilience/validator`
- `github.com/PlatformCore/Foundation/security/encryption`
### Persistence
- `github.com/PlatformCore/Foundation/persistence/tx`
- `github.com/PlatformCore/Foundation/persistence/uow`

### Messaging
- `github.com/PlatformCore/Foundation/messaging/outbox`
- `github.com/PlatformCore/Foundation/messaging/inbox`
- `github.com/PlatformCore/Foundation/messaging/dlq`
- `github.com/PlatformCore/Foundation/messaging/redrive`
- `github.com/PlatformCore/Foundation/messaging/replay`
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



