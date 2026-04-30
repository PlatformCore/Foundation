package codec

import "encoding/json"

type JSON struct{}

func (JSON) Name() string                    { return "json" }
func (JSON) Marshal(v any) ([]byte, error)   { return json.Marshal(v) }
func (JSON) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
