package core

type Handler func(*Context) error
type ErrorHandler func(*Context, error) error

func NoopHandler(*Context) error { return nil }
func HandlerFunc(fn func(*Context) error) Handler {
	if fn == nil {
		return NoopHandler
	}
	return Handler(fn)
}
