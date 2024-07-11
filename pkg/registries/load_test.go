package registries

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	actual, err := LoadConfig("./testdata/registries.json")
	require.NoError(t, err)

	expected := Config{
		Version: "1.0",
		Registries: []Registry{
			{
				Prefix: "docker.io",
				Mirrors: []string{
					"mirror1.example.com/dockerhub",
					"mirror2.example.com/dockerhub",
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

	require.Equal(t, expected, actual)
}

func TestLoadConfig_invalid(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  string
	}{
		{
			name: "unsupported-version",
			err:  "expected registry config schema version 1.0",
		},
		{
			name: "unsupported-field",
			err:  `unknown field "unsupported_field"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := fmt.Sprintf("registries-%s.json", tc.name)
			file = filepath.Join("testdata", file)
			_, err := LoadConfig(file)
			require.ErrorContains(t, err, tc.err)
		})
	}
}
