package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveCommand(t *testing.T) {
	cfgFile := filepath.Join("..", "..", "pkg", "registries", "testdata", "registries.json")
	buf := bytes.Buffer{}
	os.Args = []string{"testee", "--config", cfgFile, "resolve", "golang:latest"}
	stdout = &buf

	err := Execute()
	require.NoError(t, err)

	expectStdout := `mirror1.example.com/dockerhub/library/golang:latest
mirror2.example.com/dockerhub/library/golang:latest
docker.io/library/golang:latest
`

	require.Equal(t, expectStdout, buf.String(), "stdout")
}

func TestResolveCommand_writes_output_to_file(t *testing.T) {
	cfgFile := filepath.Join("..", "..", "pkg", "registries", "testdata", "registries.json")

	tmpDir, err := os.MkdirTemp("", "registry-conf-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	outputFile := filepath.Join(tmpDir, "imagerefs.txt")

	var buf bytes.Buffer

	os.Args = []string{"testee", "--config", cfgFile, "resolve", "golang:latest", outputFile}
	stdout = &buf
	err = Execute()
	require.NoError(t, err)

	actual, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	expected := `mirror1.example.com/dockerhub/library/golang:latest
mirror2.example.com/dockerhub/library/golang:latest
docker.io/library/golang:latest
`

	require.Equal(t, expected, string(actual), "output file")
	require.Equal(t, "", buf.String(), "stdout")

	t.Run("fail if file exists", func(t *testing.T) {
		err = Execute()
		require.ErrorContains(t, err, "file exists")
	})
}

func TestResolveCommand_invalid_arguments(t *testing.T) {
	os.Args = []string{"testee", "resolve"}
	err := Execute()
	require.ErrorContains(t, err, "invalid arguments")
}
