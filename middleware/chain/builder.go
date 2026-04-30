package chain

import "github.com/PlatformCore/libpackage/transport/core"

type Builder struct{ chain *Chain }

func NewBuilder(name string) *Builder              { return &Builder{chain: New(name)} }
func (b *Builder) Use(mw core.Middleware) *Builder { b.chain.Use(mw); return b }
func (b *Builder) UseIf(enabled bool, mw core.Middleware) *Builder {
	if enabled {
		b.chain.Use(mw)
	}
	return b
}
func (b *Builder) Error(h core.ErrorHandler) *Builder { b.chain.OnError(h); return b }
func (b *Builder) Build() *Chain                      { return b.chain }
