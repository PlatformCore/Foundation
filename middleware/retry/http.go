package retry

import "github.com/PlatformCore/libpackage/transport/core"

func HTTP(cfg Config) core.Middleware { return Middleware(cfg) }
