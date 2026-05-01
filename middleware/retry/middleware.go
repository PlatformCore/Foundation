package retry

import (
	"context"

	goretry "github.com/PlatformCore/libpackage/resilience/retry"
	"github.com/PlatformCore/libpackage/transport/core"
)

type Config = goretry.Config

func DefaultConfig() Config { return goretry.DefaultConfig() }

func Middleware(cfg Config) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			res := goretry.Do(ctx.Context, cfg, func(cctx context.Context) error { ctx.WithContext(cctx); return next(ctx) })
			return res.Err
		}
	}
}
