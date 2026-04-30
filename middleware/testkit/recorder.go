package testkit

import "github.com/PlatformCore/libpackage/transport/core"

type Recorder struct {
	Before int
	After  int
}

func (r *Recorder) Middleware() core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error { r.Before++; err := next(ctx); r.After++; return err }
	}
}
