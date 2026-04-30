package config

type MiddlewareItem struct {
	Name      string         `json:"name" yaml:"name"`
	Enabled   bool           `json:"enabled" yaml:"enabled"`
	Options   map[string]any `json:"options" yaml:"options"`
	AppliesTo []string       `json:"applies_to" yaml:"applies_to"`
}
type Config struct {
	Profile     string           `json:"profile" yaml:"profile"`
	Middlewares []MiddlewareItem `json:"middlewares" yaml:"middlewares"`
}

func Default() Config { return Config{Profile: "production"} }
