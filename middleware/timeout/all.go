package timeout

import (
	"context"
	grpcx "github.com/PlatformCore/libpackage/transport/grpc"
	thttp "github.com/PlatformCore/libpackage/transport/http"
	kafkax "github.com/PlatformCore/libpackage/transport/kafka"
	mqttx "github.com/PlatformCore/libpackage/transport/mqtt"
	natsx "github.com/PlatformCore/libpackage/transport/nats"
	redisx "github.com/PlatformCore/libpackage/transport/redis"
	"time"
)

func HTTP(d time.Duration) thttp.Middleware {
	return func(next thttp.Handler) thttp.Handler {
		return func(ctx *thttp.Context) error {
			c, cancel := context.WithTimeout(ctx.Context.Context, d)
			defer cancel()
			ctx.Context.Context = c
			return next(ctx)
		}
	}
}
func GRPCUnary(d time.Duration) grpcx.UnaryMiddleware {
	return func(next grpcx.UnaryHandler) grpcx.UnaryHandler {
		return func(ctx *grpcx.UnaryContext) (any, error) {
			c, cancel := context.WithTimeout(ctx.Context.Context, d)
			defer cancel()
			ctx.Context.Context = c
			return next(ctx)
		}
	}
}
func NATS(d time.Duration) natsx.Middleware {
	return func(next natsx.Handler) natsx.Handler {
		return func(ctx *natsx.Context) error {
			c, cancel := context.WithTimeout(ctx.Context.Context, d)
			defer cancel()
			ctx.Context.Context = c
			return next(ctx)
		}
	}
}
func Kafka(d time.Duration) kafkax.Middleware {
	return func(next kafkax.Handler) kafkax.Handler {
		return func(ctx *kafkax.Context) error {
			c, cancel := context.WithTimeout(ctx.Context.Context, d)
			defer cancel()
			ctx.Context.Context = c
			return next(ctx)
		}
	}
}
func MQTT(d time.Duration) mqttx.Middleware {
	return func(next mqttx.Handler) mqttx.Handler {
		return func(ctx *mqttx.Context) error {
			c, cancel := context.WithTimeout(ctx.Context.Context, d)
			defer cancel()
			ctx.Context.Context = c
			return next(ctx)
		}
	}
}
func Redis(d time.Duration) redisx.Middleware {
	return func(next redisx.Handler) redisx.Handler {
		return func(ctx *redisx.Context) error {
			c, cancel := context.WithTimeout(ctx.Context.Context, d)
			defer cancel()
			ctx.Context.Context = c
			return next(ctx)
		}
	}
}
