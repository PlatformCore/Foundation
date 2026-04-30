package main

import (
	adapter "github.com/PlatformCore/libpackage/adapter/http"
	"github.com/PlatformCore/libpackage/transport/core"
	"net/http"
)

func main() {
	h := func(ctx *core.Context) error { return ctx.Response.JSON(map[string]any{"ok": true}) }
	http.HandleFunc("/", adapter.NetHTTP(h))
	_ = http.ListenAndServe(":8080", nil)
}
