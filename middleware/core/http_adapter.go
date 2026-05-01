package core

import (
	"net/http"

	transportcore "github.com/PlatformCore/libpackage/transport/core"
)

// HTTP adapts the canonical transport/core middleware into net/http middleware.
// This prevents maintaining a second HTTP-only chain engine.
func HTTP(mw transportcore.Middleware, operation string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := transportcore.New(r.Context(), transportcore.TransportHTTP, operation)
			if operation == "" {
				ctx.Operation = r.Method + " " + r.URL.Path
			}
			ctx.Request.Method = r.Method
			ctx.Request.Path = r.URL.Path
			for k, vals := range r.Header {
				if len(vals) > 0 {
					ctx.Metadata.Set(k, vals[0])
					ctx.Request.Metadata.Set(k, vals[0])
				}
			}
			h := func(c *transportcore.Context) error {
				next.ServeHTTP(w, r.WithContext(c.Context))
				return nil
			}
			if mw != nil {
				h = mw(h)
			}
			if err := h(ctx); err != nil {
				e := transportcore.ToError(err)
				status := e.Status
				if status <= 0 {
					status = http.StatusInternalServerError
				}
				http.Error(w, e.Error(), status)
			}
		})
	}
}
