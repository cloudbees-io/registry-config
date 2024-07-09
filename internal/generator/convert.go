package generator

import (
	"io"

	"github.com/BurntSushi/toml"
	"github.com/calculi-corp/registry-config/pkg/registry"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
)

// Convert converts the cloudbees registry JSON configuration into the
// Redhat system registry configuration file format.
func Convert(registryConfig registry.ContainerRegistry) sysregistriesv2.V2RegistriesConf {
	var sysRegistries sysregistriesv2.V2RegistriesConf

	var registries = make([]sysregistriesv2.Registry, 0, len(registryConfig.Registries))

	for _, r := range registryConfig.Registries {
		registry := sysregistriesv2.Registry{
			Prefix: r.Prefix,
			Endpoint: sysregistriesv2.Endpoint{
				Location: r.Location,
			},
			Blocked:            false,
			MirrorByDigestOnly: false,
		}
		for _, m := range r.Mirrors {
			registry.Mirrors = append(registry.Mirrors, sysregistriesv2.Endpoint{
				Location: m,
				Insecure: false,
			})
		}

		registries = append(registries, registry)
	}

	sysRegistries.Registries = registries

	return sysRegistries
}

func Write(registryConfig sysregistriesv2.V2RegistriesConf, w io.Writer) error {
	return toml.NewEncoder(w).Encode(registryConfig)
}
