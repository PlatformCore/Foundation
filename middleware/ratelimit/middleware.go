package ratelimit

import (
	"context"
	"errors"
	"strconv"

	rl "github.com/PlatformCore/libpackage/ratelimit"
	"github.com/PlatformCore/libpackage/transport/core"
)

type KeyFunc func(*core.Context) rl.Key

type Limiter interface {
	Allow(context.Context, rl.Key) (rl.Result, error)
}

type Options struct {
	Limiter   Limiter
	Namespace string
	KeyFunc   KeyFunc
	FailOpen  bool
}

func defaultKey(namespace string) KeyFunc {
	return func(ctx *core.Context) rl.Key {
		id := ctx.UserID
		if id == "" {
			id = ctx.Metadata.Get("x-api-key")
		}
		if id == "" {
			id = ctx.RequestID
		}
		if id == "" {
			id = ctx.Operation
		}
		if namespace == "" {
			namespace = ctx.Transport
		}
		return rl.Key{Namespace: namespace, Identity: id}
	}
}

func Middleware(opts Options) core.Middleware {
	if opts.KeyFunc == nil {
		opts.KeyFunc = defaultKey(opts.Namespace)
	}
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			if opts.Limiter == nil {
				return next(ctx)
			}
			res, err := opts.Limiter.Allow(ctx.Context, opts.KeyFunc(ctx))
			if err != nil && !errors.Is(err, rl.ErrLimited) {
				if opts.FailOpen {
					return next(ctx)
				}
				return err
			}
			ctx.Response.Headers.Set("x-ratelimit-limit", strconv.FormatInt(res.Limit, 10))
			ctx.Response.Headers.Set("x-ratelimit-remaining", strconv.FormatInt(res.Remaining, 10))
			if !res.Allowed || errors.Is(err, rl.ErrLimited) {
				return core.NewError("RATE_LIMITED", "rate limit exceeded", 429)
			}
			return next(ctx)
		}
	}
}
