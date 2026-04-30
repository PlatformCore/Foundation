package bridge

import "github.com/PlatformCore/libpackage/transport/core"

type Adapter interface {
	Name() string
	Wrap(core.Handler, ...core.Middleware) any
}

func Chain(h core.Handler, mws ...core.Middleware) core.Handler { return core.Chain(h, mws...) }
