//go:build legacy
// +build legacy

package idempotency

import (
	"context"
	"fmt"
	grpcx "github.com/PlatformCore/libpackage/transport/grpc"
	thttp "github.com/PlatformCore/libpackage/transport/http"
	kafkax "github.com/PlatformCore/libpackage/transport/kafka"
	mqttx "github.com/PlatformCore/libpackage/transport/mqtt"
	natsx "github.com/PlatformCore/libpackage/transport/nats"
	redisx "github.com/PlatformCore/libpackage/transport/redis"
)

type Store interface {
	Seen(context.Context, string) (bool, error)
	Save(context.Context, string) error
}

func check(ctx context.Context, s Store, k string) error {
	if s == nil || k == "" {
		return nil
	}
	seen, err := s.Seen(ctx, k)
	if err != nil {
		return err
	}
	if seen {
		return fmt.Errorf("duplicate request: %s", k)
	}
	return s.Save(ctx, k)
}
func HTTP(s Store) thttp.Middleware {
	return func(next thttp.Handler) thttp.Handler {
		return func(ctx *thttp.Context) error {
			if err := check(ctx.Context.Context, s, ctx.Request.Header.Get("Idempotency-Key")); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}
func GRPCUnary(s Store) grpcx.UnaryMiddleware {
	return func(next grpcx.UnaryHandler) grpcx.UnaryHandler {
		return func(ctx *grpcx.UnaryContext) (any, error) {
			if err := check(ctx.Context.Context, s, ctx.RequestID); err != nil {
				return nil, err
			}
			return next(ctx)
		}
	}
}
func NATS(s Store) natsx.Middleware {
	return func(next natsx.Handler) natsx.Handler {
		return func(ctx *natsx.Context) error {
			if err := check(ctx.Context.Context, s, ctx.Message.Key); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}
func Kafka(s Store) kafkax.Middleware {
	return func(next kafkax.Handler) kafkax.Handler {
		return func(ctx *kafkax.Context) error {
			if err := check(ctx.Context.Context, s, ctx.Message.Key); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}
func MQTT(s Store) mqttx.Middleware {
	return func(next mqttx.Handler) mqttx.Handler {
		return func(ctx *mqttx.Context) error {
			if err := check(ctx.Context.Context, s, ctx.Message.Key); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}
func Redis(s Store) redisx.Middleware {
	return func(next redisx.Handler) redisx.Handler {
		return func(ctx *redisx.Context) error {
			if err := check(ctx.Context.Context, s, ctx.Message.Key); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}
