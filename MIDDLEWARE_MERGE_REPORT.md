# Middleware Merge Report

Merged source: `middleware_all.zip`
Base project: `libpackage-enterprise-complete.zip`

## Result

- Incoming zip files: 32
- Base Go files before merge: 206
- Go files after merge: 231
- Original source files were preserved. Existing files were not deleted.
- When an incoming file overlapped with a newer/fixed file, the existing file was kept and the incoming file was saved as `*_incoming.go`.

## Important decisions

1. `middleware/` from the incoming zip was imported into `middleware/nethttp/`.
2. `obslib/pkg/middleware/http` was imported into `middleware/obshttp/`.
3. `obslib/pkg/middleware/grpc` was imported into `middleware/obsgrpc/`.
4. Existing `middleware/ids`, `middleware/trace`, `middleware/propagation`, `middleware/registry`, `middleware/clock`, and `middleware/pool` were kept because the latest project already had those packages; incoming variants were preserved when different.
5. `obslib/pkg/telemetry/provider.go` was imported into `observability/telemetry/provider/`.
6. Go module files were updated for middleware and observability dependencies.

## Mapping

See `MIDDLEWARE_MERGE_MAPPING.csv` for exact source-to-target mapping.

## Build-safety note

Incoming overlap variants were stored as `*_incoming.go.txt` so they do not create duplicate Go symbols during build.

Renamed overlap variants:
- `middleware/event/middleware_incoming.go` -> `middleware/event/middleware_incoming.go.txt`
- `middleware/ids/ids_incoming.go` -> `middleware/ids/ids_incoming.go.txt`
- `middleware/propagation/propagation_incoming.go` -> `middleware/propagation/propagation_incoming.go.txt`
- `middleware/trace/span_incoming.go` -> `middleware/trace/span_incoming.go.txt`
