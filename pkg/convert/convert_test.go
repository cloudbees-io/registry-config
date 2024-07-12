package convert

import (
	"bytes"
	"testing"

	"github.com/calculi-corp/registry-config/pkg/registries"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/stretchr/testify/require"
)

func TestToRegistriesConf(t *testing.T) {
	input := registries.Config{
		Registries: []registries.Registry{
			{
				Prefix: "docker.io/library",
				Mirrors: []string{
					"mirror1.example.com/dockerhub_lib",
					"mirror2.example.com/dockerhub_lib",
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
				Prefix: "quay.io",
				Mirrors: []string{
					"mirror1.example.com/quay",
					"mirror2.example.com/quay",
				},
			},
		},
	}
	expected := sysregistriesv2.V2RegistriesConf{
		UnqualifiedSearchRegistries: []string{"docker.io"},
		ShortNameMode:               "enforcing",
		Registries: []sysregistriesv2.Registry{
			{
				Endpoint: sysregistriesv2.Endpoint{Location: "docker.io/library"},
				Mirrors: []sysregistriesv2.Endpoint{
					{Location: "mirror1.example.com/dockerhub_lib"},
					{Location: "mirror2.example.com/dockerhub_lib"},
				},
			},
			{
				Endpoint: sysregistriesv2.Endpoint{Location: "docker.io/myorg"},
				Mirrors: []sysregistriesv2.Endpoint{
					{Location: "mirror1.example.com/myorg"},
					{Location: "mirror2.example.com/myorg"},
				},
			},
			{
				Endpoint: sysregistriesv2.Endpoint{Location: "quay.io"},
				Mirrors: []sysregistriesv2.Endpoint{
					{Location: "mirror1.example.com/quay"},
					{Location: "mirror2.example.com/quay"},
				},
			},
		},
	}

	actual := ToRegistriesConf(input)

	require.Equal(t, expected, actual)
}

func TestWrite(t *testing.T) {
	conf := sysregistriesv2.V2RegistriesConf{
		Registries: []sysregistriesv2.Registry{
			{
				Endpoint: sysregistriesv2.Endpoint{Location: "docker.io"},
				Mirrors: []sysregistriesv2.Endpoint{
					{
						Location: "mirror1.example.com/dockerhub",
					},
				},
			},
		},
	}

	var buf bytes.Buffer

	err := Write(conf, &buf)
	require.NoError(t, err)

	expected := `short-name-mode = ""

[[registry]]
  prefix = ""
  location = "docker.io"

  [[registry.mirror]]
    location = "mirror1.example.com/dockerhub"
`

	require.Equal(t, expected, buf.String())
}
