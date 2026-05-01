package preset

import (
	"time"

	mididem "github.com/PlatformCore/libpackage/messaging/idempotency"
	"github.com/PlatformCore/libpackage/middleware/chain"
	"github.com/PlatformCore/libpackage/middleware/idempotency"
	"github.com/PlatformCore/libpackage/middleware/logging"
	"github.com/PlatformCore/libpackage/middleware/metrics"
	"github.com/PlatformCore/libpackage/middleware/ratelimit"
	"github.com/PlatformCore/libpackage/middleware/recovery"
	"github.com/PlatformCore/libpackage/middleware/retry"
	"github.com/PlatformCore/libpackage/middleware/timeout"
	obslog "github.com/PlatformCore/libpackage/observability/logging"
	rl "github.com/PlatformCore/libpackage/resilience/ratelimit"
	goretry "github.com/PlatformCore/libpackage/resilience/retry"
)

type EnterpriseOptions struct {
	Logger      obslog.Logger
	RateLimiter *rl.Limiter
	Idempotency *mididem.Manager
	Timeout     time.Duration
	Retry       goretry.Config
}

// UnifiedEnterprise is the recommended single middleware chain.
// It reuses the existing heavy engines instead of creating duplicate retry/logger/metrics/ratelimit logic.
func UnifiedEnterprise(opts EnterpriseOptions) *chain.Chain {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.Retry.MaxAttempts == 0 {
		opts.Retry = goretry.DefaultConfig()
	}
	c := chain.New("unified-enterprise")
	c.Use(
		recovery.Middleware(),
		logging.Middleware(&logging.Options{Logger: opts.Logger}),
		metrics.Middleware(nil),
		timeout.Middleware(opts.Timeout),
		retry.Middleware(opts.Retry),
	)
	if opts.RateLimiter != nil {
		c.Use(ratelimit.Middleware(ratelimit.Options{Limiter: opts.RateLimiter}))
	}
	if opts.Idempotency != nil {
		c.Use(idempotency.Middleware(idempotency.Options{Manager: opts.Idempotency}))
	}
	return c
}
