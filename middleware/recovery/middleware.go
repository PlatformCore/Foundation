package recovery

import (
	"fmt"
	"runtime/debug"

	"github.com/PlatformCore/libpackage/transport/core"
)

type PanicHandler func(*core.Context, any, []byte)

type Options struct {
	OnPanic      PanicHandler
	IncludeStack bool
}

func Middleware(opts ...Options) core.Middleware {
	o := Options{IncludeStack: true}
	if len(opts) > 0 {
		o = opts[0]
	}
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					var stack []byte
					if o.IncludeStack {
						stack = debug.Stack()
					}
					if o.OnPanic != nil {
						o.OnPanic(ctx, r, stack)
					}
					err = core.Wrap("PANIC", fmt.Sprint(r), 500, nil)
				}
			}()
			return next(ctx)
		}
	}
}
