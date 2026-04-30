package nats

import "github.com/PlatformCore/libpackage/transport/core"

type Message struct {
	Subject string
	Key     string
	Headers core.Metadata
	Payload []byte
	Raw     any
}
type Handler func(*Context) error
type Middleware func(Handler) Handler
type AckFunc func() error

type Context struct {
	*core.Context
	Message Message
	Ack     AckFunc
	Nak     AckFunc
}

func NewContext(operation string, msg Message) *Context {
	return &Context{Context: core.New(nil, "nats", operation), Message: msg}
}
func Chain(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] != nil {
			h = mws[i](h)
		}
	}
	return h
}
