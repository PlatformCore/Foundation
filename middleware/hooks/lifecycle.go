package hooks

import "github.com/PlatformCore/libpackage/transport/core"

type Lifecycle struct{ hooks []Hook }

func New(h ...Hook) *Lifecycle  { return &Lifecycle{hooks: h} }
func (l *Lifecycle) Use(h Hook) { l.hooks = append(l.hooks, h) }
func (l *Lifecycle) Middleware() core.Middleware {
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			for _, h := range l.hooks {
				if err := h.Before(ctx); err != nil {
					return err
				}
			}
			err := next(ctx)
			for i := len(l.hooks) - 1; i >= 0; i-- {
				err = l.hooks[i].After(ctx, err)
			}
			return err
		}
	}
}
