package tracing

import "github.com/PlatformCore/libpackage/transport/core"

func REDIS() core.Middleware { return Middleware(nil) }
