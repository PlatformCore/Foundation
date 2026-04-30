package chain

import "github.com/PlatformCore/libpackage/transport/core"

type Chain struct {
	name        string
	middlewares []core.Middleware
	onError     core.ErrorHandler
}

func New(name string, mws ...core.Middleware) *Chain { return &Chain{name: name, middlewares: mws} }
func (c *Chain) Name() string                        { return c.name }
func (c *Chain) Use(mw ...core.Middleware) *Chain {
	c.middlewares = append(c.middlewares, mw...)
	return c
}
func (c *Chain) OnError(h core.ErrorHandler) *Chain { c.onError = h; return c }
func (c *Chain) Then(h core.Handler) core.Handler {
	wrapped := core.Chain(h, c.middlewares...)
	if c.onError == nil {
		return wrapped
	}
	return func(ctx *core.Context) error {
		if err := wrapped(ctx); err != nil {
			return c.onError(ctx, err)
		}
		return nil
	}
}
