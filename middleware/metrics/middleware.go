package metrics

import (
	"time"

	gometrics "github.com/PlatformCore/libpackage/observability/metrics"
	"github.com/PlatformCore/libpackage/transport/core"
)

type Options struct{ Namespace string }

type Engine struct {
	requests *gometrics.Counter
	failures *gometrics.Counter
	latency  *gometrics.Histogram
}

func NewEngine(opts Options) *Engine {
	ns := opts.Namespace
	if ns == "" {
		ns = "libpackage_transport"
	}
	return &Engine{
		requests: gometrics.NewCounter(ns+"_requests_total", "Total transport requests", "transport", "operation", "status"),
		failures: gometrics.NewCounter(ns+"_failures_total", "Total transport failures", "transport", "operation", "code"),
		latency:  gometrics.NewHistogram(ns+"_latency_seconds", "Transport request latency", []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}, "transport", "operation", "status"),
	}
}

func Middleware(opts *Options) core.Middleware {
	var o Options
	if opts != nil {
		o = *opts
	}
	e := NewEngine(o)
	return func(next core.Handler) core.Handler {
		return func(ctx *core.Context) error {
			start := time.Now()
			err := next(ctx)
			status := "ok"
			code := "none"
			if err != nil {
				status = "error"
				code = "UNKNOWN"
				if te, ok := err.(*core.TransportError); ok {
					code = te.Code
				}
			}
			e.requests.Inc(ctx.Transport, ctx.Operation, status)
			e.latency.Observe(time.Since(start).Seconds(), ctx.Transport, ctx.Operation, status)
			if err != nil {
				e.failures.Inc(ctx.Transport, ctx.Operation, code)
			}
			return err
		}
	}
}

func Default() core.Middleware { return Middleware(nil) }
