package mqtt

import t "github.com/PlatformCore/libpackage/transport/mqtt"

type RawMessage struct {
	Subject string
	Key     string
	Headers map[string]string
	Payload []byte
	Raw     any
}

func Wrap(h t.Handler, mws ...t.Middleware) func(RawMessage) error {
	chain := t.Chain(h, mws...)
	return func(r RawMessage) error {
		msg := t.Message{Subject: r.Subject, Key: r.Key, Headers: r.Headers, Payload: r.Payload, Raw: r.Raw}
		return chain(t.NewContext(r.Subject, msg))
	}
}
