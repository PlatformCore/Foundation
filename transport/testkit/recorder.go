package testkit

import "github.com/PlatformCore/libpackage/transport/core"

type Recorder struct {
	Calls    int
	Contexts []*core.Context
	Errors   []error
}

func (r *Recorder) Handler(err error) core.Handler {
	return func(ctx *core.Context) error {
		r.Calls++
		r.Contexts = append(r.Contexts, ctx)
		if err != nil {
			r.Errors = append(r.Errors, err)
		}
		return err
	}
}
