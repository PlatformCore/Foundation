package core

import "time"

type Request struct {
	Transport  string
	Operation  string
	Method     string
	Path       string
	Subject    string
	Topic      string
	Key        string
	Metadata   Metadata
	Body       []byte
	Raw        any
	ReceivedAt time.Time
}

func NewRequest(transport, operation string) *Request {
	return &Request{Transport: transport, Operation: operation, Metadata: Metadata{}, ReceivedAt: time.Now()}
}
