package ratelimit

import "github.com/PlatformCore/libpackage/transport/core"

func Redis(opts Options) core.Middleware { return Middleware(opts) }
