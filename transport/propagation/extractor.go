package propagation

import "github.com/PlatformCore/libpackage/transport/core"

type Extractor struct{}

func (Extractor) Extract(c Carrier) core.Metadata {
	md := core.Metadata{}
	if c == nil {
		return md
	}
	for _, k := range c.Keys() {
		md.Set(k, c.Get(k))
	}
	return md
}
func ApplyToContext(ctx *core.Context, md core.Metadata) {
	if ctx == nil {
		return
	}
	for k, v := range md {
		ctx.WithMetadata(k, v)
	}
	ctx.RequestID = md.Get(core.HeaderRequestID)
	ctx.TraceID = md.Get(core.HeaderTraceID)
	ctx.UserID = md.Get(core.HeaderUserID)
	ctx.TenantID = md.Get(core.HeaderTenantID)
}
