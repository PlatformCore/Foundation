package http

import (
	"github.com/PlatformCore/libpackage/transport/core"
	"io"
	"net/http"
)

func NetHTTP(h core.Handler, mws ...core.Middleware) http.HandlerFunc {
	wrapped := core.Chain(h, mws...)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := core.New(r.Context(), core.TransportHTTP, r.Method+" "+r.URL.Path)
		ctx.Request.Method = r.Method
		ctx.Request.Path = r.URL.Path
		ctx.Request.Raw = r
		for k, v := range r.Header {
			if len(v) > 0 {
				ctx.Request.Metadata.Set(k, v[0])
			}
		}
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			ctx.Request.Body = b
		}
		err := wrapped(ctx)
		for k, v := range ctx.Response.Metadata {
			w.Header().Set(k, v)
		}
		if err != nil {
			e := core.ToError(err)
			w.WriteHeader(e.Status)
			_, _ = w.Write([]byte(e.Error()))
			return
		}
		w.WriteHeader(ctx.Response.StatusCode)
		_, _ = w.Write(ctx.Response.Body)
	}
}
