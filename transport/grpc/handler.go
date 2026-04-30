package grpc

import "github.com/PlatformCore/libpackage/transport/core"

type UnaryHandler func(*UnaryContext) (any, error)
type UnaryMiddleware func(UnaryHandler) UnaryHandler
type StreamHandler func(*StreamContext) error
type StreamMiddleware func(StreamHandler) StreamHandler

type UnaryContext struct {
	*core.Context
	FullMethod string
	Request    any
}
type StreamContext struct {
	*core.Context
	FullMethod string
	Stream     any
}

func NewUnary(operation string, req any) *UnaryContext {
	return &UnaryContext{Context: core.New(nil, "grpc", operation), FullMethod: operation, Request: req}
}
func ChainUnary(h UnaryHandler, mws ...UnaryMiddleware) UnaryHandler {
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			h = mws[i](h)
		}
	}
	return h
}
func ChainStream(h StreamHandler, mws ...StreamMiddleware) StreamHandler {
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			h = mws[i](h)
		}
	}
	return h
}
