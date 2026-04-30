package tracing

import "github.com/PlatformCore/libpackage/transport/core"

func HTTP() core.Middleware { return Middleware(nil) }
