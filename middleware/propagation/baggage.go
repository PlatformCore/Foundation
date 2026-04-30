package propagation

import "github.com/PlatformCore/libpackage/transport/core"

func Baggage(keys ...string) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			for _, k := range keys {
				if v := ctx.Metadata.Get(k); v != "" {
					ctx.Set(k, v)
				}
			}
			return next(ctx)
		}
	}
}
