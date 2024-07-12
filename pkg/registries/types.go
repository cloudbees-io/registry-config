package registries

// Config defines the container registry mirror mapping.
type Config struct {
	Version    string     `json:"version"`
	Registries []Registry `json:"registries"`
}

// Registry defines the mirrors that should be used for a registry.
type Registry struct {
	Prefix  string   `json:"prefix"`
	Mirrors []string `json:"mirrors"`
}
