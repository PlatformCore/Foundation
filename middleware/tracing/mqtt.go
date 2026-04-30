package tracing

import "github.com/PlatformCore/libpackage/transport/core"

func MQTT() core.Middleware { return Middleware(nil) }
