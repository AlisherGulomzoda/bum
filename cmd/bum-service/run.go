//nolint:forbidigo // here we haven't loaded our own logger yet so that we use log lib
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bum-service/internal/app"

	"github.com/spf13/cobra"
)

// RunCmd is run cmd command.
func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Run the service",
		Long:    "backend server",
		PreRunE: loadConfigs,
		RunE:    runServer,
	}

	cmd.Flags().String(configPathFlag, defaultConfigPath, "path to config yml file")

	return cmd
}

const (
	configPathFlag    = "config_path"
	defaultConfigPath = "config.yml"
)

// runServer runs the Server.
func runServer(_ *cobra.Command, _ []string) error {
	log.Println("Starting the program ...")

	application, err := app.NewService(cfg)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	defer application.Close()

	err = application.Run()
	if err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	<-waitSignalToExit()

	return nil
}

// waitSignalToExit creates a signal channel that signals the application to exit.
func waitSignalToExit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return ch
}
