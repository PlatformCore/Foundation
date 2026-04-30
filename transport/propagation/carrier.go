package propagation

type Carrier interface {
	Get(string) string
	Set(string, string)
	Keys() []string
}

type MapCarrier map[string]string

func (m MapCarrier) Get(k string) string { return m[k] }
func (m MapCarrier) Set(k, v string)     { m[k] = v }
func (m MapCarrier) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
