package ratelimit

import "github.com/PlatformCore/libpackage/transport/core"

func Kafka(opts Options) core.Middleware { return Middleware(opts) }
