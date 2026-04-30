package testkit

import "github.com/PlatformCore/libpackage/transport/core"

func NewRequest(transport, operation string) *core.Context {
	return core.New(nil, transport, operation)
}
