package grpcadapter

import "github.com/PlatformCore/libpackage/transport/core"

type UnaryHandler func(ctx any, req any) (any, error)

func Unary(h core.Handler, mws ...core.Middleware) UnaryHandler {
	wrapped := core.Chain(h, mws...)
	return func(rawCtx any, req any) (any, error) {
		ctx := core.New(nil, core.TransportGRPC, "grpc.unary")
		ctx.Request.Raw = req
		err := wrapped(ctx)
		return ctx.Result.Data, err
	}
}
