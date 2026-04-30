package nethttp

import (
	"net/http"
	thttp "github.com/PlatformCore/libpackage/transport/http"
)

func Handler(h thttp.Handler, onError func(http.ResponseWriter,*http.Request,error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := thttp.NewContext(w,r,r.Method+" "+r.URL.Path)
		if err := h(ctx); err != nil && onError != nil { onError(w,r,err) }
	}
}
func Chain(h thttp.Handler, mws ...thttp.Middleware) http.HandlerFunc { return Handler(thttp.Chain(h,mws...), nil) }
