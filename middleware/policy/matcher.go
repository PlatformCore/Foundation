package policy

import "github.com/PlatformCore/libpackage/transport/core"

type Matcher func(*core.Context) bool

func Transport(name string) Matcher {
	return func(ctx *core.Context) bool { return ctx != nil && ctx.Transport == name }
}
func Operation(op string) Matcher {
	return func(ctx *core.Context) bool { return ctx != nil && ctx.Operation == op }
}
func And(ms ...Matcher) Matcher {
	return func(ctx *core.Context) bool {
		for _, m := range ms {
			if m != nil && !m(ctx) {
				return false
			}
		}
		return true
	}
}
