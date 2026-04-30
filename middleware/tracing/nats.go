package tracing

import "github.com/PlatformCore/libpackage/transport/core"

func NATS() core.Middleware { return Middleware(nil) }
