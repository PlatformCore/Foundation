package policy

import "github.com/PlatformCore/libpackage/transport/core"

type Rule struct {
	Match  Matcher
	Use    []string
	Allow  bool
	Reason string
}
type Rules []Rule

func (rs Rules) Decide(ctx *core.Context) Decision {
	for _, r := range rs {
		if r.Match == nil || r.Match(ctx) {
			return Decision{Allow: r.Allow || len(r.Use) > 0, Reason: r.Reason, MiddlewareNames: r.Use}
		}
	}
	return Decision{Allow: true}
}
