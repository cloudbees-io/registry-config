package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cloudbees-io/registry-config/pkg/convert"
	"github.com/cloudbees-io/registry-config/pkg/registries"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converted the registry config into Red Hat's registries.conf format",
	Long: "Converts the given registry configuration into Red Hat's registries.conf " +
		"file format and writes it to the specified path.",
	RunE: convertRegistriesConfig,
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

func convertRegistriesConfig(_ *cobra.Command, args []string) (err error) {
	if len(args) != 1 {
		return fmt.Errorf("invalid arguments: %d arguments given but requires 1 argument: OUTFILE", len(args))
	}

	config, err := registries.LoadConfig(configFile)
	if err != nil {
		return err
	}

	rhConfig := convert.ToRegistriesConf(config)

	outFile, err := os.OpenFile(args[0], os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0640)
	if err != nil {
		return fmt.Errorf("write registries.conf: %w", err)
	}

	defer func() {
		e := outFile.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	err = convert.Write(rhConfig, outFile)
	if err != nil {
		return fmt.Errorf("write registries.conf: %w", err)
	}

	return err
}
