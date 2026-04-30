package core

import "strings"

type Metadata map[string]string

func (m Metadata) Clone() Metadata {
	out := Metadata{}
	for k, v := range m {
		out[k] = v
	}
	return out
}
func (m Metadata) Set(k, v string) {
	if k == "" {
		return
	}
	m[CanonicalKey(k)] = v
}
func (m Metadata) Get(k string) string { return m[CanonicalKey(k)] }
func (m Metadata) Del(k string)        { delete(m, CanonicalKey(k)) }
func CanonicalKey(k string) string     { return strings.ToLower(strings.TrimSpace(k)) }

const (
	HeaderRequestID      = "x-request-id"
	HeaderTraceID        = "traceparent"
	HeaderUserID         = "x-user-id"
	HeaderTenantID       = "x-tenant-id"
	HeaderIdempotencyKey = "idempotency-key"
)
