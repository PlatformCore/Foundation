# Separate Module Installation Guide

Updated: 2026-04-28

Install any module independently:
```bash
go get github.com/driftappdev/<module-path>@<version>
```

Examples:
```bash
go get github.com/PlatformCore/Foundation/platform/eventbus@v0.1.0
go get github.com/PlatformCore/Foundation/platform/servicemesh@v0.1.0
go get github.com/PlatformCore/Foundation/ratelimit@v1.0.0
go get github.com/PlatformCore/Foundation/security@v1.0.0
```

Tag format used in this repository:
```
<module-path>/vX.Y.Z
```

