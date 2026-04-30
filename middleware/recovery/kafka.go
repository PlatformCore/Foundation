package recovery

import "github.com/PlatformCore/libpackage/transport/core"

func KAFKA() core.Middleware { return Middleware(nil) }
