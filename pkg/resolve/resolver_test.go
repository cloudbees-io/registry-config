package resolve

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/calculi-corp/registry-config/pkg/registries"
)

func TestResolve(t *testing.T) {
	config := registries.Config{
		Registries: []registries.Registry{
			{
				Prefix: "docker.io/library",
				Mirrors: []string{
					"mirror1.example.com/dockerhub",
					"mirror2.example.com/dockerhub",
				},
			},
			{
				Prefix: "docker.io/myorg",
				Mirrors: []string{
					"mirror1.example.com/myorg",
					"mirror2.example.com/myorg",
				},
			},
			{
				Prefix: "docker.io/otherorg",
				Mirrors: []string{
					"mirror1.example.com/extorg",
					"mirror2.example.com/extorg",
				},
			},
			{
				Prefix: "quay.io",
				Mirrors: []string{
					"mirror1.example.com/quay",
					"mirror2.example.com/quay",
				},
			},
		},
	}

	resolver, err := NewResolver(config)
	require.NoError(t, err)
	defer resolver.Close()

	tests := []struct {
		name     string
		imageRef string
		expected []string
	}{
		{
			name:     "match dockerhub library mirror",
			imageRef: "golang:1.21",
			expected: []string{
				"mirror1.example.com/dockerhub/golang:1.21",
				"mirror2.example.com/dockerhub/golang:1.21",
				"docker.io/library/golang:1.21",
			},
		},
		{
			name:     "match dockerhub org mirror",
			imageRef: "myorg/myimage:1.2.3",
			expected: []string{
				"mirror1.example.com/myorg/myimage:1.2.3",
				"mirror2.example.com/myorg/myimage:1.2.3",
				"docker.io/myorg/myimage:1.2.3",
			},
		},
		{
			name:     "match dockerhub org mirror for other org",
			imageRef: "otherorg/myimage:1.2.3",
			expected: []string{
				"mirror1.example.com/extorg/myimage:1.2.3",
				"mirror2.example.com/extorg/myimage:1.2.3",
				"docker.io/otherorg/myimage:1.2.3",
			},
		},
		{
			name:     "match dockerhub org mirror for absolute image ref",
			imageRef: "docker.io/myorg/myimage:1.2.3",
			expected: []string{
				"mirror1.example.com/myorg/myimage:1.2.3",
				"mirror2.example.com/myorg/myimage:1.2.3",
				"docker.io/myorg/myimage:1.2.3",
			},
		},
		{
			name:     "match quay.io",
			imageRef: "quay.io/podman/stable:5.1.1",
			expected: []string{
				"mirror1.example.com/quay/podman/stable:5.1.1",
				"mirror2.example.com/quay/podman/stable:5.1.1",
				"quay.io/podman/stable:5.1.1",
			},
		},
		{
			name:     "return original location for non-matching image ref",
			imageRef: "registry.example.com/someimage:1.2.3",
			expected: []string{
				"registry.example.com/someimage:1.2.3",
			},
		},
		{
			name:     "match with digest support",
			imageRef: "golang:1.21@sha256:2a03a6059f21e150ae84b0973863609494aad70f0a80eaeb64bddd8d92465812",
			expected: []string{
				"mirror1.example.com/dockerhub/golang@sha256:2a03a6059f21e150ae84b0973863609494aad70f0a80eaeb64bddd8d92465812",
				"mirror2.example.com/dockerhub/golang@sha256:2a03a6059f21e150ae84b0973863609494aad70f0a80eaeb64bddd8d92465812",
				"docker.io/library/golang@sha256:2a03a6059f21e150ae84b0973863609494aad70f0a80eaeb64bddd8d92465812",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := resolver.Resolve(tc.imageRef)
			require.NoError(t, err)
			require.Equalf(t, tc.expected, out, "imageref: %s", tc.imageRef)
		})
	}

	err = resolver.Close()
	require.NoError(t, err, "Close()")
}
