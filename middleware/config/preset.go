package config

func ProductionPreset() Config {
	return Config{Profile: "production", Middlewares: []MiddlewareItem{{Name: "request_id", Enabled: true}, {Name: "recovery", Enabled: true}, {Name: "logging", Enabled: true}, {Name: "tracing", Enabled: true}, {Name: "metrics", Enabled: true}, {Name: "ratelimit", Enabled: true}}}
}
func DevelopmentPreset() Config {
	return Config{Profile: "development", Middlewares: []MiddlewareItem{{Name: "request_id", Enabled: true}, {Name: "recovery", Enabled: true}, {Name: "logging", Enabled: true}}}
}
