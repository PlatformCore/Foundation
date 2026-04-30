package idempotency

import "github.com/PlatformCore/libpackage/transport/core"

func GRPC(opts Options) core.Middleware { return Middleware(opts) }
