package propagation

import "github.com/PlatformCore/libpackage/transport/core"

func Trace() core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			if ctx.TraceID == "" {
				ctx.TraceID = ctx.Metadata.Get(core.HeaderTraceID)
			}
			return next(ctx)
		}
	}
}
