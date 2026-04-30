package idempotency

import "github.com/PlatformCore/libpackage/transport/core"

func NATS(opts Options) core.Middleware { return Middleware(opts) }
