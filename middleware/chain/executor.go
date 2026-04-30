package chain

import "github.com/PlatformCore/libpackage/transport/core"

type Executor struct {
	chain   *Chain
	handler core.Handler
}

func NewExecutor(c *Chain, h core.Handler) *Executor { return &Executor{chain: c, handler: c.Then(h)} }
func (e *Executor) Execute(ctx *core.Context) error  { return e.handler(ctx) }
