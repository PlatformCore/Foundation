package core

type Middleware func(Handler) Handler

func Chain(h Handler, mws ...Middleware) Handler {
	if h == nil {
		h = NoopHandler
	}
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			h = mws[i](h)
		}
	}
	return h
}
func Compose(mws ...Middleware) Middleware {
	return func(next Handler) Handler { return Chain(next, mws...) }
}
func Named(name string, mw Middleware) Middleware { return mw }
