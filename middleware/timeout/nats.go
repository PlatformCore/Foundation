package timeout

import (
	"github.com/PlatformCore/libpackage/transport/core"
	"time"
)

func NATS(d time.Duration) core.Middleware { return Middleware(d) }
