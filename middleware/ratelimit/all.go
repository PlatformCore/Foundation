//go:build legacy
// +build legacy

package ratelimit

import (
	"fmt"
	grpcx "github.com/PlatformCore/libpackage/transport/grpc"
	thttp "github.com/PlatformCore/libpackage/transport/http"
	kafkax "github.com/PlatformCore/libpackage/transport/kafka"
	mqttx "github.com/PlatformCore/libpackage/transport/mqtt"
	natsx "github.com/PlatformCore/libpackage/transport/nats"
	redisx "github.com/PlatformCore/libpackage/transport/redis"
)

type Limiter interface{ Allow(key string) bool }

func key(v string) string {
	if v == "" {
		return "anonymous"
	}
	return v
}
func HTTP(l Limiter) thttp.Middleware {
	return func(next thttp.Handler) thttp.Handler {
		return func(ctx *thttp.Context) error {
			if l != nil && !l.Allow(key(ctx.UserID)) {
				return fmt.Errorf("rate limit exceeded")
			}
			return next(ctx)
		}
	}
}
func GRPCUnary(l Limiter) grpcx.UnaryMiddleware {
	return func(next grpcx.UnaryHandler) grpcx.UnaryHandler {
		return func(ctx *grpcx.UnaryContext) (any, error) {
			if l != nil && !l.Allow(key(ctx.UserID)) {
				return nil, fmt.Errorf("rate limit exceeded")
			}
			return next(ctx)
		}
	}
}
func NATS(l Limiter) natsx.Middleware {
	return func(next natsx.Handler) natsx.Handler {
		return func(ctx *natsx.Context) error {
			if l != nil && !l.Allow(key(ctx.Message.Subject)) {
				return fmt.Errorf("rate limit exceeded")
			}
			return next(ctx)
		}
	}
}
func Kafka(l Limiter) kafkax.Middleware {
	return func(next kafkax.Handler) kafkax.Handler {
		return func(ctx *kafkax.Context) error {
			if l != nil && !l.Allow(key(ctx.Message.Subject)) {
				return fmt.Errorf("rate limit exceeded")
			}
			return next(ctx)
		}
	}
}
func MQTT(l Limiter) mqttx.Middleware {
	return func(next mqttx.Handler) mqttx.Handler {
		return func(ctx *mqttx.Context) error {
			if l != nil && !l.Allow(key(ctx.Message.Subject)) {
				return fmt.Errorf("rate limit exceeded")
			}
			return next(ctx)
		}
	}
}
func Redis(l Limiter) redisx.Middleware {
	return func(next redisx.Handler) redisx.Handler {
		return func(ctx *redisx.Context) error {
			if l != nil && !l.Allow(key(ctx.Message.Subject)) {
				return fmt.Errorf("rate limit exceeded")
			}
			return next(ctx)
		}
	}
}
