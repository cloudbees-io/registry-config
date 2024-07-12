package convert

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/calculi-corp/registry-config/pkg/registries"

	"github.com/containers/image/v5/pkg/sysregistriesv2"
)

// ToRegistriesConf converts the registry config into the Red Hat's format.
func ToRegistriesConf(config registries.Config) sysregistriesv2.V2RegistriesConf {
	var registries = make([]sysregistriesv2.Registry, len(config.Registries))

	for i, r := range config.Registries {
		mirrors := make([]sysregistriesv2.Endpoint, len(r.Mirrors))
		for j, m := range r.Mirrors {
			mirrors[j] = sysregistriesv2.Endpoint{
				Location: m,
				Insecure: false,
			}
		}

		registries[i] = sysregistriesv2.Registry{
			Endpoint:           sysregistriesv2.Endpoint{Location: r.Prefix},
			Blocked:            false,
			MirrorByDigestOnly: false,
			Mirrors:            mirrors,
		}
	}

	return sysregistriesv2.V2RegistriesConf{
		UnqualifiedSearchRegistries: []string{"docker.io"},
		ShortNameMode:               "enforcing",
		Registries:                  registries,
	}
}

// Write writes a given Red Hat registries.conf into the given writer.
func Write(config sysregistriesv2.V2RegistriesConf, w io.Writer) error {
	err := toml.NewEncoder(w).Encode(config)
	if err != nil {
		return fmt.Errorf("write rh registries.conf file: %w", err)
	}

	return nil
}
