package registry

type ContainerRegistry struct {
	Version    string     `json:"version"`
	Registries []Registry `json:"registries"`
}

type Registry struct {
	Prefix   string   `json:"prefix"`
	Location string   `json:"location"`
	Mirrors  []string `json:"mirrors"`
}
