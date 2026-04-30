package grpcadapter

import (
	"context"
	grpcx "github.com/PlatformCore/libpackage/transport/grpc"
	gogrpc "google.golang.org/grpc"
)

func UnaryInterceptor(mws ...grpcx.UnaryMiddleware) gogrpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *gogrpc.UnaryServerInfo, handler gogrpc.UnaryHandler) (any, error) {
		h := func(c *grpcx.UnaryContext) (any, error) { return handler(ctx, req) }
		return grpcx.ChainUnary(h, mws...)(&grpcx.UnaryContext{FullMethod: info.FullMethod, Request: req})
	}
}
