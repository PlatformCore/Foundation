package ratelimit

import "github.com/PlatformCore/libpackage/transport/core"

func MQTT(opts Options) core.Middleware { return Middleware(opts) }
