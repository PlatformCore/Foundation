package grpcadapter

import "github.com/PlatformCore/libpackage/transport/core"

func Interceptor(h core.Handler, mws ...core.Middleware) any { return Unary(h, mws...) }
