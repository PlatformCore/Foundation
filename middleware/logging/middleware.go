package logging

import (
	"github.com/PlatformCore/libpackage/observability/logging"
	"github.com/PlatformCore/libpackage/transport/core"
)

type Engine = logging.Logger

type Options struct {
	Logger          Engine
	IncludeMetadata bool
}

func Middleware(opts *Options) core.Middleware {
	if opts == nil {
		opts = &Options{}
	}
	lg := opts.Logger
	if lg == nil {
		lg = logging.New()
	}
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			fields := []any{"transport", ctx.Transport, "operation", ctx.Operation, "request_id", ctx.RequestID, "trace_id", ctx.TraceID, "user_id", ctx.UserID, "tenant_id", ctx.TenantID}
			lg.InfoContext(ctx.Context, "request started", fields...)
			err := next(ctx)
			fields = append(fields, "duration_ms", ctx.Duration().Milliseconds())
			if err != nil {
				lg.ErrorContext(ctx.Context, "request failed", append(fields, "error", err.Error())...)
				return err
			}
			lg.InfoContext(ctx.Context, "request completed", fields...)
			return nil
		}
	}
}

func Default() core.Middleware { return Middleware(nil) }
