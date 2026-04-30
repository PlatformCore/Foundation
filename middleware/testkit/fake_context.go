package testkit

import "github.com/PlatformCore/libpackage/transport/core"

func FakeContext() *core.Context { return core.New(nil, core.TransportHTTP, "test") }
