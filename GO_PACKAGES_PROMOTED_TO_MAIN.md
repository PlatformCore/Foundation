# Go* Packages Promoted To Main

This reorganization promotes the stronger `go*` implementations from compatibility/legacy subfolders into first-class enterprise package locations.

Rules used:
- No `.go` source file was deleted.
- Basic/lightweight implementations remain in their original/basic package groups.
- `go*` implementations are now the preferred import paths.
- Old compatibility locations are left as empty/deprecated folders only when needed for documentation; source logic is moved to the main package path.

## Main package mapping

| Package | Old path | New main path | Basic/secondary remains |
|---|---|---|---|
| goerror | `core/errors/goerror_compat` | `core/goerror` | `core/errors` |
| gologger | `observability/logging/gologger` | `observability/gologger` | `observability/logging` |
| gometrics | `observability/gometrics` | `observability/gometrics` | `observability/metrics/basic` |
| gotracing | `observability/tracing/gotracing` | `observability/gotracing` | `observability/tracing` |
| goratelimit | `ratelimit/compat/goratelimit` | `ratelimit/goratelimit` | `ratelimit/{limiter,policy,key,memory_store,redis_store,...}` |
| goauth | `security/auth_middleware/goauth` | `security/goauth` | `security/auth_middleware` |
| gosanitizer | `security/sanitizer/gosanitizer` | `security/gosanitizer` | `security/sanitizer` |
| gocircuit | `resilience/circuitbreaker/gocircuit` | `resilience/gocircuit` | `resilience/circuitbreaker` |
| goretry | `resilience/retry/goretry` | `resilience/goretry` | `resilience/retry` |
| gotimeout | `resilience/timeout/gotimeout` | `resilience/gotimeout` | `resilience/timeout` |

## Preferred imports

```go
import "github.com/PlatformCore/Foundation/observability/gometrics"
import "github.com/PlatformCore/Foundation/observability/gologger"
import "github.com/PlatformCore/Foundation/observability/gotracing"
import "github.com/PlatformCore/Foundation/core/goerror"
import "github.com/PlatformCore/Foundation/security/goauth"
import "github.com/PlatformCore/Foundation/security/gosanitizer"
import "github.com/PlatformCore/Foundation/ratelimit/goratelimit"
import "github.com/PlatformCore/Foundation/resilience/gocircuit"
import "github.com/PlatformCore/Foundation/resilience/goretry"
import "github.com/PlatformCore/Foundation/resilience/gotimeout"
```

## Verification

- Go files before promotion: 227
- Go files after promotion: 227
- Result: no Go source file loss if counts match.

## Notes

- Old source import paths were searched in Go files; no direct references were found.
- Top-level enterprise modules remain unchanged: core, observability, security, ratelimit, resilience, etc.
- Basic implementations are kept as secondary packages.
