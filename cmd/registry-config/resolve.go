package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cloudbees-io/registry-config/pkg/registries"
	"github.com/cloudbees-io/registry-config/pkg/resolve"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolves a given image reference to a list of locations",
	Long:  "Resolves a given image reference to a list of locations",
	RunE:  resolveImageReference,
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}

func resolveImageReference(_ *cobra.Command, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("invalid arguments: %d arguments given but requires 1-2 arguments: IMAGEREF [OUTFILE]", len(args))
	}

	config, err := registries.LoadConfig(configFile)
	if err != nil {
		return err
	}

	resolver, err := resolve.NewResolver(config)
	if err != nil {
		return err
	}
	defer resolver.Close()

	locations, err := resolver.Resolve(args[0])
	if err != nil {
		return err
	}

	output := strings.Join(locations, "\n")

	if len(args) == 2 && args[1] != "-" {
		err := writeToFile(args[1], []byte(output+"\n"))
		if err != nil {
			return fmt.Errorf("write registries.conf: %w", err)
		}

		return nil
	}

	fmt.Fprintln(stdout, output)

	return nil
}

func writeToFile(path string, content []byte) (err error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}

	defer func() {
		e := file.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	_, err = io.Copy(file, bytes.NewReader(content))

	return err
}
