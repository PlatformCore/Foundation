# Foundation Release Notes (EN) - Dev + PM

Release date: 2026-04-29  
Repository: https://github.com/PlatformCore/Foundation

## Release Summary
- Aligned all Foundation module paths to `github.com/PlatformCore/Foundation/...`
- Kept a multi-module layout for independent installation and release
- Added centralized version manifest at `tools/module_versions.json`
- Added one-command release flow (bump + commit + push + tag)

## Module Versions (This Release)
| Module | Path | Version | Status | Best Use |
|---|---|---|---|---|
| clients | `clients` | v0.1.0 | stable | HTTP/gRPC/NATS client layer |
| config | `config` | v0.1.0 | stable | config loading/env/file/watcher |
| core | `core` | v0.1.0 | stable | foundational contracts/types/errors/logger |
| messaging | `messaging` | v0.1.0 | stable | inbox/outbox/replay/dlq |
| middleware | `middleware` | v0.1.0 | stable | shared HTTP/GRPC middleware |
| observability | `observability` | v0.1.0 | stable | logging/metrics/tracing/profiler |
| orchestration | `orchestration` | v0.1.0 | stable | DI/workflow/orchestration |
| persistence | `persistence` | v0.1.0 | stable | uow/tx/distributed lock |
| platform | `platform` | v0.1.0 | stable | eventbus/featureflag/servicemesh |
| plugins | `plugins` | v0.1.0 | stable | plugin engine/registry/hooks |
| ratelimit | `ratelimit` | v0.1.0 | stable | rate limiting policy/store |
| resilience | `resilience` | v0.1.0 | stable | retry/circuit/timeout/bulkhead |
| runtime | `runtime` | v0.1.0 | stable | lifecycle/shutdown/signals |
| security | `security` | v0.1.0 | stable | auth/jwt/oauth2/encryption |
| tools | `tools` | v0.1.0 | beta | release/generation helper tooling |
| validation | `validation` | v0.1.0 | stable | binding/schema validation |

## One-Command Version Bump + Release
Use:
`./scripts/release_one_command.ps1 -ModulePath platform -Version v0.2.0`

The script will:
1. Update `tools/module_versions.json`
2. Commit to `main`
3. Push `main` to `origin`
4. Create and push tag in `modulePath/version` format, e.g. `platform/v0.2.0`

## PM Notes
- Each module can be released independently for phased rollouts.
- Keep `v0.x` during hardening; promote to `v1.0.0` once APIs/behavior are stable.
