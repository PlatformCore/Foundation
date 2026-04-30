module github.com/PlatformCore/libpackage/adapter

go 1.23.0

require (
	github.com/PlatformCore/libpackage/transport v0.0.0
	google.golang.org/grpc v1.80.0
)

replace github.com/PlatformCore/libpackage/transport => ../transport
