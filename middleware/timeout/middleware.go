package timeout

import (
	"context"
	"time"

	gotimeout "github.com/PlatformCore/libpackage/resilience/timeout/gotimeout"
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
			return gotimeout.Do(ctx.Context, opts.Duration, func(cctx context.Context) error { ctx.WithContext(cctx); return next(ctx) })
		}
	}
}
