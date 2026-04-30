package main

import (
	adapter "github.com/PlatformCore/libpackage/adapter/nats"
	"github.com/PlatformCore/libpackage/transport/core"
)

func main() {
	_ = adapter.Core(func(ctx *core.Context) error { ctx.Result = core.Success("ok"); return nil })
}
