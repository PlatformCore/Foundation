// Package middleware provides enterprise-grade HTTP middleware for Go services.
// It offers a composable, zero-dependency-on-framework middleware chain
// compatible with net/http and any router that supports http.Handler.
//
// Usage:
//
//	chain := middleware.New(
//	    middleware.WithRequestID(),
//	    middleware.WithLogger(logger),
//	    middleware.WithRecovery(),
//	    middleware.WithMetrics(registry),
//	    middleware.WithTracing("my-service"),
//	    middleware.WithCORS(corsConfig),
//	    middleware.WithAuth(authConfig),
//	    middleware.WithRateLimit(rateLimitConfig),
//	    middleware.WithTimeout(30*time.Second),
//	    middleware.WithCircuitBreaker(cbConfig),
//	    middleware.WithCache(cacheConfig),
//	)
//
//	http.Handle("/", chain.Then(myHandler))
package httpmw

import "net/http"

// Middleware defines a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain holds an ordered list of middleware.
type Chain struct {
	middlewares []Middleware
}

// New creates a new middleware Chain with the provided middlewares applied
// in the order given (first middleware is outermost).
func New(mws ...Middleware) *Chain {
	c := &Chain{}
	c.middlewares = make([]Middleware, len(mws))
	copy(c.middlewares, mws)
	return c
}

// Append returns a new Chain with additional middlewares appended.
func (c *Chain) Append(mws ...Middleware) *Chain {
	newMws := make([]Middleware, len(c.middlewares)+len(mws))
	copy(newMws, c.middlewares)
	copy(newMws[len(c.middlewares):], mws)
	return &Chain{middlewares: newMws}
}

// Prepend returns a new Chain with additional middlewares prepended.
func (c *Chain) Prepend(mws ...Middleware) *Chain {
	newMws := make([]Middleware, len(mws)+len(c.middlewares))
	copy(newMws, mws)
	copy(newMws[len(mws):], c.middlewares)
	return &Chain{middlewares: newMws}
}

// Then wraps the given handler with all middlewares in the chain.
// The first middleware in the chain is the outermost (first to receive request).
func (c *Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}

// ThenFunc wraps the given handler function with all middlewares in the chain.
func (c *Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}
