package registries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfig loads the registry mirror configuration file.
func LoadConfig(file string) (Config, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return Config{}, fmt.Errorf("load registries config: %w", err)
	}

	m := map[string]interface{}{}

	err = json.Unmarshal(raw, &m)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshal registries config at %s: %w", file, err)
	}

	if m["version"] != "1.0" {
		return Config{}, fmt.Errorf("expected registry config schema version 1.0 within %s but was %q", file, m["version"])
	}

	config := Config{}
	d := json.NewDecoder(bytes.NewReader(raw))

	d.DisallowUnknownFields()

	err = d.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshal registries config at %s: %w", file, err)
	}

	return config, nil
}
