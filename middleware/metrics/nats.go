package metrics

import "github.com/PlatformCore/libpackage/transport/core"

func NATS(opts ...*Options) core.Middleware {
	if len(opts) > 0 {
		return Middleware(opts[0])
	}
	return Middleware(nil)
}
