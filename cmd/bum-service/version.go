package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

// VersionCmd is version command.
func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print binary version",
		Long:  "Show binary version (build number).",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = fmt.Fprintf(
				cmd.OutOrStdout(),
				"%s version %s\n",
				cmd.Root().Use,
				version,
			)
		},
	}

	cmd.Flags().String(configPathFlag, defaultConfigPath, "path to config yml file")

	return cmd
}
