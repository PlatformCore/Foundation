package http

import (
	"github.com/PlatformCore/libpackage/transport/core"
	"net/http"
)

type Handler func(*Context) error
type Middleware func(Handler) Handler

type Context struct {
	*core.Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func NewContext(w http.ResponseWriter, r *http.Request, operation string) *Context {
	return &Context{Context: core.New(r.Context(), "http", operation), Request: r, ResponseWriter: w}
}
func Chain(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			h = mws[i](h)
		}
	}
	return h
}
