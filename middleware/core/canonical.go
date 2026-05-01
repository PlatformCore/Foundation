package core

import (
	"time"

	mwlogging "github.com/PlatformCore/libpackage/middleware/logging"
	mwmetrics "github.com/PlatformCore/libpackage/middleware/metrics"
	mwratelimit "github.com/PlatformCore/libpackage/middleware/ratelimit"
	mwrecovery "github.com/PlatformCore/libpackage/middleware/recovery"
	mwretry "github.com/PlatformCore/libpackage/middleware/retry"
	mwtimeout "github.com/PlatformCore/libpackage/middleware/timeout"
	mwtracing "github.com/PlatformCore/libpackage/middleware/tracing"
	transportcore "github.com/PlatformCore/libpackage/transport/core"
)

// EnterpriseOptions wires canonical wrappers only once. It keeps advanced logic
// in resilience and observability packages but makes the middleware layer simple.
type EnterpriseOptions struct {
	Recovery  *mwrecovery.Options
	Logging   *mwlogging.Options
	Metrics   *mwmetrics.Options
	Tracing   mwtracing.Tracer
	Retry     *mwretry.Config
	Timeout   time.Duration
	RateLimit *mwratelimit.Options
}

// EnterpriseRegistry returns a canonical ordered registry. Add custom business
// middleware after this registry; do not add retry/timeout/ratelimit again.
func EnterpriseRegistry(opts EnterpriseOptions) *Registry {
	reg := NewRegistry()
	if opts.Recovery != nil {
		reg.Register(Spec{Name: "recovery", Stage: StageIngress, Priority: 10, MW: mwrecovery.Middleware(*opts.Recovery)})
	} else {
		reg.Register(Spec{Name: "recovery", Stage: StageIngress, Priority: 10, MW: mwrecovery.Middleware()})
	}
	reg.Register(Spec{Name: "logging", Stage: StageObservability, Priority: 10, MW: mwlogging.Middleware(opts.Logging)})
	reg.Register(Spec{Name: "metrics", Stage: StageObservability, Priority: 20, MW: mwmetrics.Middleware(opts.Metrics)})
	if opts.Tracing != nil {
		reg.Register(Spec{Name: "tracing", Stage: StageObservability, Priority: 30, MW: mwtracing.Middleware(opts.Tracing)})
	}
	if opts.Retry != nil {
		reg.Register(Spec{Name: "retry", Stage: StageResilience, Priority: 10, MW: mwretry.Middleware(*opts.Retry)})
	}
	if opts.Timeout > 0 {
		reg.Register(Spec{Name: "timeout", Stage: StageResilience, Priority: 20, MW: mwtimeout.Middleware(opts.Timeout)})
	}
	if opts.RateLimit != nil {
		reg.Register(Spec{Name: "ratelimit", Stage: StagePolicy, Priority: 10, MW: mwratelimit.Middleware(*opts.RateLimit)})
	}
	return reg
}

func Enterprise(opts EnterpriseOptions) transportcore.Middleware {
	return transportcore.Compose(EnterpriseRegistry(opts).Middleware()...)
}
