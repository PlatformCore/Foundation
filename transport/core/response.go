package core

import "encoding/json"

type Response struct {
	StatusCode int
	Metadata   Metadata
	Body       []byte
	Raw        any
}

func NewResponse() *Response { return &Response{StatusCode: 200, Metadata: Metadata{}} }
func (r *Response) Header(k, v string) *Response {
	if r.Metadata == nil {
		r.Metadata = Metadata{}
	}
	r.Metadata.Set(k, v)
	return r
}
func (r *Response) Status(code int) *Response { r.StatusCode = code; return r }
func (r *Response) JSON(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.Body = b
	r.Header("content-type", "application/json")
	r.Raw = v
	return nil
}
func (r *Response) Bytes(b []byte) *Response { r.Body = b; return r }
