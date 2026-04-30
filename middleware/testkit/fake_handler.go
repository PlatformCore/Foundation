package testkit

import "github.com/PlatformCore/libpackage/transport/core"

func FakeHandler(err error) core.Handler { return func(*core.Context) error { return err } }
