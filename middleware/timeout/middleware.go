package timeout

import (
	"context"
	"time"

	timeoutcore "github.com/PlatformCore/libpackage/resilience/timeout"
	"github.com/PlatformCore/libpackage/transport/core"
)

type Options struct{ Duration time.Duration }

func Middleware(d time.Duration) core.Middleware { return WithOptions(Options{Duration: d}) }
func WithOptions(opts Options) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			if opts.Duration <= 0 {
				return next(ctx)
			}
			return timeout.Do(ctx.Context, opts.Duration, func(cctx context.Context) error { ctx.WithContext(cctx); return next(ctx) })
		}
	}
}
