package timeout

import (
	"github.com/PlatformCore/libpackage/transport/core"
	"time"
)

func MQTT(d time.Duration) core.Middleware { return Middleware(d) }
