package registry

import "github.com/PlatformCore/libpackage/transport/core"

func Static(mw core.Middleware) Factory {
	return func(Config) (core.Middleware, error) { return mw, nil }
}
func MustRegister(name string, mw core.Middleware) { Default.Register(name, Static(mw)) }
