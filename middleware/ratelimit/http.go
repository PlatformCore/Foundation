package ratelimit

import "github.com/PlatformCore/libpackage/transport/core"

func HTTP(opts Options) core.Middleware { return Middleware(opts) }
