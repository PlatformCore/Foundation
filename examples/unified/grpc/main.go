package main

import (
	adapter "github.com/PlatformCore/libpackage/adapter/grpc"
	"github.com/PlatformCore/libpackage/transport/core"
)

func main() {
	_ = adapter.Unary(func(ctx *core.Context) error { ctx.Result = core.Success("ok"); return nil })
}
