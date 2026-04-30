package mqtt

import "github.com/PlatformCore/libpackage/transport/core"

type Message struct {
	Key     string
	Subject string
	Topic   string
	Headers map[string]string
	Payload []byte
	Raw     any
}

func Publisher(h core.Handler, mws ...core.Middleware) func(Message) error {
	wrapped := core.Chain(h, mws...)
	return func(msg Message) error {
		ctx := core.New(nil, core.TransportMQTT, "mqtt.publisher")
		ctx.Request.Key = msg.Key
		ctx.Request.Subject = msg.Subject
		ctx.Request.Topic = msg.Topic
		ctx.Request.Body = msg.Payload
		ctx.Request.Raw = msg.Raw
		for k, v := range msg.Headers {
			ctx.Request.Metadata.Set(k, v)
		}
		return wrapped(ctx)
	}
}
