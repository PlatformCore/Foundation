package idempotency

import (
	midem "github.com/PlatformCore/libpackage/messaging/idempotency"
	"github.com/PlatformCore/libpackage/transport/core"
)

type KeyFunc func(*core.Context) string

type Options struct {
	Manager     *midem.Manager
	Namespace   string
	ServiceName string
	KeyFunc     KeyFunc
}

func Middleware(opts Options) core.Middleware {
	if opts.KeyFunc == nil {
		opts.KeyFunc = func(ctx *core.Context) string { return ctx.Metadata.Get(core.HeaderIdempotencyKey) }
	}
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			if opts.Manager == nil {
				return next(ctx)
			}
			key := opts.KeyFunc(ctx)
			if key == "" {
				return next(ctx)
			}
			ns := opts.Namespace
			if ns == "" {
				ns = ctx.Transport
			}
			res, err := opts.Manager.Evaluate(ctx.Context, &midem.Request{Key: key, Namespace: ns, Payload: ctx.Request.Body, UserID: ctx.UserID, ServiceName: opts.ServiceName, OperationName: ctx.Operation, Metadata: ctx.Metadata})
			if err != nil {
				return err
			}
			if res != nil && res.IsCached {
				ctx.Result = core.Success(map[string]any{"idempotency": "cached", "request_id": res.RequestID})
				return nil
			}
			err = next(ctx)
			if err != nil {
				_ = opts.Manager.Fail(ctx.Context, ns, key, err.Error())
				return err
			}
			return opts.Manager.Complete(ctx.Context, ns, key, ctx.Result.Data, ctx.Response.StatusCode)
		}
	}
}
