package propagation

import "github.com/PlatformCore/libpackage/transport/core"

type Injector struct{}

func (Injector) Inject(ctx *core.Context, c Carrier) {
	if ctx == nil || c == nil {
		return
	}
	for k, v := range ctx.Metadata {
		c.Set(k, v)
	}
	if ctx.RequestID != "" {
		c.Set(core.HeaderRequestID, ctx.RequestID)
	}
	if ctx.TraceID != "" {
		c.Set(core.HeaderTraceID, ctx.TraceID)
	}
	if ctx.UserID != "" {
		c.Set(core.HeaderUserID, ctx.UserID)
	}
	if ctx.TenantID != "" {
		c.Set(core.HeaderTenantID, ctx.TenantID)
	}
}
