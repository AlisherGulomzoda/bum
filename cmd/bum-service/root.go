//nolint:forbidigo // here we haven't configured our logger yet, so we can use log instead
package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"bum-service/cmd/bum-service/populate"
	"bum-service/config"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bum-service",
		Long:    "backend service of bum education platform",
		PreRunE: loadConfigs,
		RunE:    runServer,
		Version: "v0.0.0",
	}

	cmd.Flags().String(configPathFlag, defaultConfigPath, "path to config yml file")

	cmd.AddCommand(RunCmd())
	cmd.AddCommand(MigrateCmd())
	cmd.AddCommand(VersionCmd())
	cmd.AddCommand(populate.Cmd())

	return cmd
}

func loadConfigs(cmd *cobra.Command, _ []string) (err error) {
	log.Println("Starting the program ...")

	log.Println("Loading configs ...")

	configPath, err := cmd.Flags().GetString(configPathFlag)
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}
