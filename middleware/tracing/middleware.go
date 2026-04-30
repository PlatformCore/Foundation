package tracing

import "github.com/PlatformCore/libpackage/transport/core"

type Tracer interface {
	Start(*core.Context) (func(error), error)
}

func Middleware(t Tracer) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			var end func(error)
			if t != nil {
				e, err := t.Start(ctx)
				if err != nil {
					return err
				}
				end = e
			}
			err := next(ctx)
			if end != nil {
				end(err)
			}
			return err
		}
	}
}
