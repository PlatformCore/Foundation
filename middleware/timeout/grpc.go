package timeout

import (
	"github.com/PlatformCore/libpackage/transport/core"
	"time"
)

func GRPC(d time.Duration) core.Middleware { return Middleware(d) }
