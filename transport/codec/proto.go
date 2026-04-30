package codec

import "errors"

type Proto struct{}

func (Proto) Name() string { return "proto" }
func (Proto) Marshal(v any) ([]byte, error) {
	m, ok := v.(interface{ Marshal() ([]byte, error) })
	if ok {
		return m.Marshal()
	}
	return nil, errors.New("proto codec requires value with Marshal()")
}
func (Proto) Unmarshal(b []byte, v any) error {
	m, ok := v.(interface{ Unmarshal([]byte) error })
	if ok {
		return m.Unmarshal(b)
	}
	return errors.New("proto codec requires value with Unmarshal([]byte)")
}
