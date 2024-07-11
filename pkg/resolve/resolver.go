package resolve

import (
	"fmt"
	"os"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/containers/image/v5/types"

	"github.com/calculi-corp/registry-config/pkg/convert"
	"github.com/calculi-corp/registry-config/pkg/registries"
)

// Resolver resolves image references using a given registry mirror config.
type Resolver struct {
	config *types.SystemContext
}

// Close closes the resolver.
func (r *Resolver) Close() error {
	if r.config != nil {
		err := os.Remove(r.config.SystemRegistriesConfPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		r.config = nil
	}

	return nil
}

// NewResolver creates a new image reference resolver using the provided config.
func NewResolver(config registries.Config) (*Resolver, error) {
	tmpRhConfFile, err := createTempRegistriesConf(config)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		config: &types.SystemContext{
			SystemRegistriesConfPath: tmpRhConfFile,
		},
	}, nil
}

// Resolve returns a list of (mirror) locations for the given image reference.
func (r *Resolver) Resolve(ref string) ([]string, error) {
	parsed, err := reference.ParseDockerRef(ref)
	if err != nil {
		return nil, err
	}

	refStr := imageRefString(parsed)

	registry, err := sysregistriesv2.FindRegistry(r.config, refStr)
	if err != nil {
		return nil, fmt.Errorf("find registry for image ref %q: %w", refStr, err)
	}

	if registry == nil {
		return []string{refStr}, nil
	}

	pullSources, err := registry.PullSourcesFromReference(parsed)
	if err != nil {
		return nil, fmt.Errorf("resolve mirrors for image ref %q: %w", refStr, err)
	}

	resolvedRefs := make([]string, len(pullSources))
	for i, src := range pullSources {
		resolvedRefs[i] = imageRefString(src.Reference)
	}

	return resolvedRefs, nil
}

func imageRefString(ref reference.Named) string {
	imgRef := ref.Name()
	if nt, ok := ref.(reference.NamedTagged); ok {
		imgRef = fmt.Sprintf("%s:%s", imgRef, nt.Tag())
	}

	if cref, ok := ref.(reference.Canonical); ok {
		imgRef = fmt.Sprintf("%s@%s", imgRef, cref.Digest())
	}

	return imgRef
}

func createTempRegistriesConf(config registries.Config) (string, error) {
	rhRegistriesConf := convert.ToRegistriesConf(config)

	tmpFile, err := os.CreateTemp("", "registries-conf-")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if err = convert.Write(rhRegistriesConf, tmpFile); err != nil {
		_ = os.Remove(tmpFile.Name())

		return "", err
	}

	return tmpFile.Name(), nil
}
