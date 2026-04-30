//go:build legacy
// +build legacy

package recovery

import (
	"fmt"
	grpcx "github.com/PlatformCore/libpackage/transport/grpc"
	thttp "github.com/PlatformCore/libpackage/transport/http"
	kafkax "github.com/PlatformCore/libpackage/transport/kafka"
	mqttx "github.com/PlatformCore/libpackage/transport/mqtt"
	natsx "github.com/PlatformCore/libpackage/transport/nats"
	redisx "github.com/PlatformCore/libpackage/transport/redis"
)

type RecoverFunc func(any) error

func recoverErr(fn RecoverFunc) error {
	if r := recover(); r != nil {
		if fn != nil {
			return fn(r)
		}
		return fmt.Errorf("panic recovered: %v", r)
	}
	return nil
}
func HTTP(fn RecoverFunc) thttp.Middleware {
	return func(next thttp.Handler) thttp.Handler {
		return func(ctx *thttp.Context) (err error) {
			defer func() {
				if e := recoverErr(fn); e != nil {
					err = e
				}
			}()
			return next(ctx)
		}
	}
}
func GRPCUnary(fn RecoverFunc) grpcx.UnaryMiddleware {
	return func(next grpcx.UnaryHandler) grpcx.UnaryHandler {
		return func(ctx *grpcx.UnaryContext) (resp any, err error) {
			defer func() {
				if e := recoverErr(fn); e != nil {
					err = e
				}
			}()
			return next(ctx)
		}
	}
}
func NATS(fn RecoverFunc) natsx.Middleware {
	return func(next natsx.Handler) natsx.Handler {
		return func(ctx *natsx.Context) (err error) {
			defer func() {
				if e := recoverErr(fn); e != nil {
					err = e
				}
			}()
			return next(ctx)
		}
	}
}
func Kafka(fn RecoverFunc) kafkax.Middleware {
	return func(next kafkax.Handler) kafkax.Handler {
		return func(ctx *kafkax.Context) (err error) {
			defer func() {
				if e := recoverErr(fn); e != nil {
					err = e
				}
			}()
			return next(ctx)
		}
	}
}
func MQTT(fn RecoverFunc) mqttx.Middleware {
	return func(next mqttx.Handler) mqttx.Handler {
		return func(ctx *mqttx.Context) (err error) {
			defer func() {
				if e := recoverErr(fn); e != nil {
					err = e
				}
			}()
			return next(ctx)
		}
	}
}
func Redis(fn RecoverFunc) redisx.Middleware {
	return func(next redisx.Handler) redisx.Handler {
		return func(ctx *redisx.Context) (err error) {
			defer func() {
				if e := recoverErr(fn); e != nil {
					err = e
				}
			}()
			return next(ctx)
		}
	}
}
