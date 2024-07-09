package resolver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	resolver, err := NewResolver("testdata/registries.json", "")
	require.NoError(t, err)

	tests := []struct {
		name     string
		imageRef string
		expected string
	}{
		{
			name:     "image and tag",
			imageRef: "golang:1.21",
			expected: "mirror1.example.com/dockerhub_lib/golang:1.21",
		},
		{
			name:     "no match",
			imageRef: "yourorg/golang:1.21",
			expected: "yourorg/golang:1.21",
		},
		{
			name:     "match",
			imageRef: "myorg/golang:1.21",
			expected: "mirror1.example.com/dockerhub/golang:1.21",
		},
		{
			name:     "quay.io match",
			imageRef: "quay.io/myorg/golang:1.21",
			expected: "mirror1.example.com/quay/golang:1.21",
		},
		{
			name:     "quay.io no match",
			imageRef: "quay.io/yourorg/golang:1.21",
			expected: "quay.io/yourorg/golang:1.21",
		},
	}

	for _, test := range tests {
		out, err := resolver.Resolve(test.imageRef)
		require.NoError(t, err)
		require.Equal(t, test.expected, out, test.name)
	}
}
