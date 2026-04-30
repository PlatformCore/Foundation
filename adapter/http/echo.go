package http

import "github.com/PlatformCore/libpackage/transport/core"

func EchoLike(h core.Handler, mws ...core.Middleware) func(any) error {
	wrapped := core.Chain(h, mws...)
	return func(raw any) error {
		ctx := core.New(nil, core.TransportHTTP, "echo")
		ctx.Request.Raw = raw
		return wrapped(ctx)
	}
}
