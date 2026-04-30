package grpcadapter

import "github.com/PlatformCore/libpackage/transport/core"

func Stream(h core.Handler, mws ...core.Middleware) func(any) error {
	wrapped := core.Chain(h, mws...)
	return func(stream any) error {
		ctx := core.New(nil, core.TransportGRPC, "grpc.stream")
		ctx.Request.Raw = stream
		return wrapped(ctx)
	}
}
