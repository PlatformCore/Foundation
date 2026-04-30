package recovery

import "github.com/PlatformCore/libpackage/transport/core"

func GRPC() core.Middleware { return Middleware(nil) }
