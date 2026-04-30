package propagation

import "github.com/PlatformCore/libpackage/transport/core"

func Middleware() core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			if ctx.Request != nil {
				for k, v := range ctx.Request.Metadata {
					ctx.WithMetadata(k, v)
				}
			}
			return next(ctx)
		}
	}
}
