package retry

import "github.com/PlatformCore/libpackage/transport/core"

func Redis(cfg Config) core.Middleware { return Middleware(cfg) }
