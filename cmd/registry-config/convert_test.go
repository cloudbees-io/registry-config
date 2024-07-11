package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertCommand(t *testing.T) {
	cfgFile := filepath.Join("..", "..", "pkg", "registries", "testdata", "registries.json")
	tmpDir, err := os.MkdirTemp("", "registry-conf-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	outputFile := filepath.Join(tmpDir, "registries.conf")

	os.Args = []string{"testee", "--config", cfgFile, "convert", outputFile}
	err = Execute()
	require.NoError(t, err)

	b, err := os.ReadFile(outputFile)
	require.NoError(t, err, "should have written output file")
	require.Contains(t, string(b), "[[registry]]\n")

	t.Run("fail if file exists", func(t *testing.T) {
		err = Execute()
		require.ErrorContains(t, err, "file exists")
	})
}

func TestConvertCommand_invalid_arguments(t *testing.T) {
	os.Args = []string{"testee", "convert"}
	err := Execute()
	require.ErrorContains(t, err, "invalid arguments")
}
