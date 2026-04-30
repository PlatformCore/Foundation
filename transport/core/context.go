package core

import (
	"context"
	"sync"
	"time"
)

const (
	TransportHTTP      = "http"
	TransportHTTP3     = "http3"
	TransportGRPC      = "grpc"
	TransportNATS      = "nats"
	TransportNATSCore  = "nats_core"
	TransportJetStream = "nats_jetstream"
	TransportKafka     = "kafka"
	TransportMQTT      = "mqtt"
	TransportRedis     = "redis"
)

type Context struct {
	Context   context.Context
	Transport string
	Operation string
	RequestID string
	TraceID   string
	UserID    string
	TenantID  string
	Metadata  Metadata
	Request   *Request
	Response  *Response
	Result    Result
	StartedAt time.Time
	Values    map[string]any
	mu        sync.RWMutex
}

func New(ctx context.Context, transport, operation string) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Context{Context: ctx, Transport: transport, Operation: operation, Metadata: Metadata{}, Request: &Request{Transport: transport, Operation: operation, Metadata: Metadata{}}, Response: NewResponse(), StartedAt: time.Now(), Values: map[string]any{}}
}

func (c *Context) Deadline() (time.Time, bool) { return c.Context.Deadline() }
func (c *Context) Done() <-chan struct{}       { return c.Context.Done() }
func (c *Context) Err() error                  { return c.Context.Err() }
func (c *Context) Value(key any) any {
	if v := c.Context.Value(key); v != nil {
		return v
	}
	if s, ok := key.(string); ok {
		return c.Get(s)
	}
	return nil
}
func (c *Context) WithContext(ctx context.Context) *Context {
	if ctx != nil {
		c.Context = ctx
	}
	return c
}
func (c *Context) Set(k string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Values == nil {
		c.Values = map[string]any{}
	}
	c.Values[k] = v
}
func (c *Context) Get(k string) any { c.mu.RLock(); defer c.mu.RUnlock(); return c.Values[k] }
func (c *Context) WithMetadata(k, v string) *Context {
	if c.Metadata == nil {
		c.Metadata = Metadata{}
	}
	c.Metadata.Set(k, v)
	if c.Request != nil {
		c.Request.Metadata.Set(k, v)
	}
	return c
}
func (c *Context) Duration() time.Duration {
	if c.StartedAt.IsZero() {
		return 0
	}
	return time.Since(c.StartedAt)
}
