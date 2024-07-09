package resolver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/calculi-corp/registry-config/internal/generator"
	"github.com/calculi-corp/registry-config/pkg/registry"
	imageparser "github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/containers/image/v5/types"
)

// Resolver is an interface for resolving an image reference using registry mirrors.
type Resolver interface {
	Resolve(ref string) (string, error)
}

type resolver struct {
	registryConfig         registry.ContainerRegistry
	registryConfigWriteDir string
}

// NewResolver creates a new Resolver using the provided registry configuration JSON.
// It also takes a base temp directory to create a directory and write the system registries configuration file in, if
// one is not passed in, it will create it's own temp directory.
func NewResolver(registryConfigFile string, registryConfigWriteBaseDir string) (Resolver, error) {
	b, err := os.ReadFile(registryConfigFile)
	if err != nil {
		return nil, err
	}

	var registryConfig registry.ContainerRegistry
	if err = json.Unmarshal(b, &registryConfig); err != nil {
		return nil, err
	}

	var registryConfigWriteDir string
	if registryConfigWriteBaseDir == "" {
		registryConfigWriteDir, err = os.MkdirTemp("", "sysregistries")
		if err != nil {
			return nil, err
		}
	} else {
		registryConfigWriteDir = filepath.Join(registryConfigWriteBaseDir, ".cloudbees", "sysregistries")
		// if directory does not exist, create it
		if _, err := os.Stat(registryConfigWriteDir); os.IsNotExist(err) {
			if err := os.MkdirAll(registryConfigWriteDir, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	return &resolver{
		registryConfig:         registryConfig,
		registryConfigWriteDir: registryConfigWriteDir}, nil
}

func (r *resolver) Resolve(ref string) (string, error) {
	parsed, err := imageparser.ParseDockerRef(ref)
	if err != nil {
		return "", err
	}

	sysRegistriesConf := generator.Convert(r.registryConfig)

	registriesConfFile := r.registryConfigWriteDir + "/registries.conf"

	f, err := os.Create(registriesConfFile)
	if err != nil {
		return "", err
	}

	if err = generator.Write(sysRegistriesConf, f); err != nil {
		return "", err
	}

	sys := &types.SystemContext{
		SystemRegistriesConfPath: registriesConfFile,
	}

	registry, err := sysregistriesv2.FindRegistry(sys, parsed.Name())
	if err != nil {
		return "", err
	}

	if registry == nil {
		return ref, nil
	}

	pullSources, err := registry.PullSourcesFromReference(parsed)
	if err != nil {
		return "", err
	}

	if len(pullSources) == 0 {
		return ref, nil
	}

	var tag string
	if nt, ok := parsed.(imageparser.NamedTagged); ok {
		tag = nt.Tag()
	}

	if len(tag) > 0 {
		return fmt.Sprintf("%s:%s", pullSources[0].Reference.Name(), tag), nil
	}

	return pullSources[0].Reference.Name(), nil
}
