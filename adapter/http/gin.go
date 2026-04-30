package http

import "github.com/PlatformCore/libpackage/transport/core"

type GinContext interface{ RequestContext() any }

func GinLike(h core.Handler, mws ...core.Middleware) func(any) error {
	wrapped := core.Chain(h, mws...)
	return func(raw any) error {
		ctx := core.New(nil, core.TransportHTTP, "gin")
		ctx.Request.Raw = raw
		return wrapped(ctx)
	}
}
