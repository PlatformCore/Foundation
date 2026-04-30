package codec

type Codec interface {
	Name() string
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}

type Registry struct{ codecs map[string]Codec }

func NewRegistry() *Registry { return &Registry{codecs: map[string]Codec{}} }
func (r *Registry) Register(c Codec) {
	if c != nil {
		r.codecs[c.Name()] = c
	}
}
func (r *Registry) Get(name string) (Codec, bool) { c, ok := r.codecs[name]; return c, ok }
