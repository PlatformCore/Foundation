package hooks

import "github.com/PlatformCore/libpackage/transport/core"

type Hook interface {
	Before(*core.Context) error
	After(*core.Context, error) error
}
type Func struct {
	BeforeFunc func(*core.Context) error
	AfterFunc  func(*core.Context, error) error
}

func (f Func) Before(ctx *core.Context) error {
	if f.BeforeFunc != nil {
		return f.BeforeFunc(ctx)
	}
	return nil
}
func (f Func) After(ctx *core.Context, err error) error {
	if f.AfterFunc != nil {
		return f.AfterFunc(ctx, err)
	}
	return err
}
func Middleware(h Hook) core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			if h != nil {
				if err := h.Before(ctx); err != nil {
					return err
				}
			}
			err := next(ctx)
			if h != nil {
				return h.After(ctx, err)
			}
			return err
		}
	}
}
