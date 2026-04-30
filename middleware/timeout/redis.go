package timeout

import (
	"github.com/PlatformCore/libpackage/transport/core"
	"time"
)

func Redis(d time.Duration) core.Middleware { return Middleware(d) }
