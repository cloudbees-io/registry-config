package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/calculi-corp/registry-config/pkg/registry"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	var registriesJson = "testdata/registries.json"

	b, err := os.ReadFile(registriesJson)
	require.NoError(t, err)

	var registryConfig registry.ContainerRegistry
	err = json.Unmarshal(b, &registryConfig)
	require.NoError(t, err)

	actualRegistries := Convert(registryConfig)

	var registriesConfigFile = "testdata/registries.conf"

	expectedRegistriesConfig, err := os.ReadFile(registriesConfigFile)
	require.NoError(t, err)

	var expectedSysRegistries sysregistriesv2.V2RegistriesConf
	require.NoError(t, toml.Unmarshal(expectedRegistriesConfig, &expectedSysRegistries))

	if !assert.Equal(t, expectedSysRegistries, actualRegistries) {
		fmt.Println("Expected registries.conf")
		Write(actualRegistries, os.Stdout)
		t.FailNow()
	}
}
