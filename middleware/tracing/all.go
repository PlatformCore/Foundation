package tracing

import (
	"context"
	grpcx "github.com/PlatformCore/libpackage/transport/grpc"
	thttp "github.com/PlatformCore/libpackage/transport/http"
	kafkax "github.com/PlatformCore/libpackage/transport/kafka"
	mqttx "github.com/PlatformCore/libpackage/transport/mqtt"
	natsx "github.com/PlatformCore/libpackage/transport/nats"
	redisx "github.com/PlatformCore/libpackage/transport/redis"
)

type Tracer interface {
	Start(context.Context, string) (context.Context, func(error))
}

func start(t Tracer, ctx context.Context, op string) (context.Context, func(error)) {
	if t == nil {
		return ctx, func(error) {}
	}
	return t.Start(ctx, op)
}
func HTTP(t Tracer) thttp.Middleware {
	return func(next thttp.Handler) thttp.Handler {
		return func(ctx *thttp.Context) error {
			c, end := start(t, ctx.Context.Context, ctx.Operation)
			ctx.Context.Context = c
			err := next(ctx)
			end(err)
			return err
		}
	}
}
func GRPCUnary(t Tracer) grpcx.UnaryMiddleware {
	return func(next grpcx.UnaryHandler) grpcx.UnaryHandler {
		return func(ctx *grpcx.UnaryContext) (any, error) {
			c, end := start(t, ctx.Context.Context, ctx.FullMethod)
			ctx.Context.Context = c
			resp, err := next(ctx)
			end(err)
			return resp, err
		}
	}
}
func NATS(t Tracer) natsx.Middleware {
	return func(next natsx.Handler) natsx.Handler {
		return func(ctx *natsx.Context) error {
			c, end := start(t, ctx.Context.Context, ctx.Operation)
			ctx.Context.Context = c
			err := next(ctx)
			end(err)
			return err
		}
	}
}
func Kafka(t Tracer) kafkax.Middleware {
	return func(next kafkax.Handler) kafkax.Handler {
		return func(ctx *kafkax.Context) error {
			c, end := start(t, ctx.Context.Context, ctx.Operation)
			ctx.Context.Context = c
			err := next(ctx)
			end(err)
			return err
		}
	}
}
func MQTT(t Tracer) mqttx.Middleware {
	return func(next mqttx.Handler) mqttx.Handler {
		return func(ctx *mqttx.Context) error {
			c, end := start(t, ctx.Context.Context, ctx.Operation)
			ctx.Context.Context = c
			err := next(ctx)
			end(err)
			return err
		}
	}
}
func Redis(t Tracer) redisx.Middleware {
	return func(next redisx.Handler) redisx.Handler {
		return func(ctx *redisx.Context) error {
			c, end := start(t, ctx.Context.Context, ctx.Operation)
			ctx.Context.Context = c
			err := next(ctx)
			end(err)
			return err
		}
	}
}
