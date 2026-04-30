package policy

import "github.com/PlatformCore/libpackage/transport/core"

type Decision struct {
	Allow           bool
	Reason          string
	MiddlewareNames []string
}
type Policy interface{ Decide(*core.Context) Decision }
type Func func(*core.Context) Decision

func (f Func) Decide(ctx *core.Context) Decision { return f(ctx) }
