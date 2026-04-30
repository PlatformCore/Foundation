package codec

type Raw struct{}

func (Raw) Name() string { return "raw" }
func (Raw) Marshal(v any) ([]byte, error) {
	if b, ok := v.([]byte); ok {
		return b, nil
	}
	if s, ok := v.(string); ok {
		return []byte(s), nil
	}
	return nil, nil
}
func (Raw) Unmarshal(b []byte, v any) error {
	if p, ok := v.(*[]byte); ok {
		*p = append((*p)[:0], b...)
	}
	return nil
}
