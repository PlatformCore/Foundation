package retry

import "github.com/PlatformCore/libpackage/transport/core"

func GRPC(cfg Config) core.Middleware { return Middleware(cfg) }
