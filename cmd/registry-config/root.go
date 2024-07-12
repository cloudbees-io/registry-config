package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var (
	rootCmd = &cobra.Command{
		Use:           "cloudbees-registry-config",
		Short:         "CLI for resolving image references",
		Long:          "CLI for resolving container image references",
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	configFile           = "/etc/cloudbees/registries.json"
	stdout     io.Writer = os.Stdout // for testing
)

func init() {
	rootCmd.SetFlagErrorFunc(handleError)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", configFile, "Path to the registries.json file")
}

func handleError(cmd *cobra.Command, err error) error {
	_ = cmd.Help()

	return err
}

func Execute() error {
	return rootCmd.Execute()
}
