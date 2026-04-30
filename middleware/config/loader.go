package config

import (
	"encoding/json"
	"io"
)

func LoadJSON(r io.Reader) (Config, error) {
	var c Config
	err := json.NewDecoder(r).Decode(&c)
	return c, err
}
