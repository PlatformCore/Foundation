module github.com/PlatformCore/libpackage/messaging/example_integration

go 1.25.0

replace github.com/PlatformCore/libpackage/messaging/dlq => ../dlq

replace github.com/PlatformCore/libpackage/messaging/inbox => ../inbox

replace github.com/PlatformCore/libpackage/messaging/outbox => ../outbox

replace github.com/PlatformCore/libpackage/messaging/redrive => ../redrive

require (
	github.com/PlatformCore/libpackage/messaging/dlq v0.0.0-00010101000000-000000000000
	github.com/PlatformCore/libpackage/messaging/inbox v0.0.0-00010101000000-000000000000
	github.com/PlatformCore/libpackage/messaging/outbox v0.0.0-00010101000000-000000000000
	github.com/PlatformCore/libpackage/messaging/redrive v0.0.0-00010101000000-000000000000
)

