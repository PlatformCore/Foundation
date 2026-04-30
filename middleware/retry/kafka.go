package retry

import "github.com/PlatformCore/libpackage/transport/core"

func Kafka(cfg Config) core.Middleware { return Middleware(cfg) }
