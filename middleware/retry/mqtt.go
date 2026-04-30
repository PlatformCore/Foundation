package retry

import "github.com/PlatformCore/libpackage/transport/core"

func MQTT(cfg Config) core.Middleware { return Middleware(cfg) }
