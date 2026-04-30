package retry

import "github.com/PlatformCore/libpackage/transport/core"

func NATS(cfg Config) core.Middleware { return Middleware(cfg) }
