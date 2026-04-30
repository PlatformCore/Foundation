//go:build legacy
// +build legacy

package metrics

import (
	grpcx "github.com/PlatformCore/libpackage/transport/grpc"
	thttp "github.com/PlatformCore/libpackage/transport/http"
	kafkax "github.com/PlatformCore/libpackage/transport/kafka"
	mqttx "github.com/PlatformCore/libpackage/transport/mqtt"
	natsx "github.com/PlatformCore/libpackage/transport/nats"
	redisx "github.com/PlatformCore/libpackage/transport/redis"
	"time"
)

type Recorder interface {
	Observe(transport, operation string, dur time.Duration, err error)
}

func observe(r Recorder, transport, op string, start time.Time, err error) {
	if r != nil {
		r.Observe(transport, op, time.Since(start), err)
	}
}
func HTTP(r Recorder) thttp.Middleware {
	return func(next thttp.Handler) thttp.Handler {
		return func(ctx *thttp.Context) error {
			st := time.Now()
			err := next(ctx)
			observe(r, "http", ctx.Operation, st, err)
			return err
		}
	}
}
func GRPCUnary(r Recorder) grpcx.UnaryMiddleware {
	return func(next grpcx.UnaryHandler) grpcx.UnaryHandler {
		return func(ctx *grpcx.UnaryContext) (any, error) {
			st := time.Now()
			resp, err := next(ctx)
			observe(r, "grpc", ctx.FullMethod, st, err)
			return resp, err
		}
	}
}
func NATS(r Recorder) natsx.Middleware {
	return func(next natsx.Handler) natsx.Handler {
		return func(ctx *natsx.Context) error {
			st := time.Now()
			err := next(ctx)
			observe(r, "nats", ctx.Operation, st, err)
			return err
		}
	}
}
func Kafka(r Recorder) kafkax.Middleware {
	return func(next kafkax.Handler) kafkax.Handler {
		return func(ctx *kafkax.Context) error {
			st := time.Now()
			err := next(ctx)
			observe(r, "kafka", ctx.Operation, st, err)
			return err
		}
	}
}
func MQTT(r Recorder) mqttx.Middleware {
	return func(next mqttx.Handler) mqttx.Handler {
		return func(ctx *mqttx.Context) error {
			st := time.Now()
			err := next(ctx)
			observe(r, "mqtt", ctx.Operation, st, err)
			return err
		}
	}
}
func Redis(r Recorder) redisx.Middleware {
	return func(next redisx.Handler) redisx.Handler {
		return func(ctx *redisx.Context) error {
			st := time.Now()
			err := next(ctx)
			observe(r, "redis", ctx.Operation, st, err)
			return err
		}
	}
}
