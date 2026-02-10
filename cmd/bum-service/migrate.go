package main

import (
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	bum "bum-service"
)

// MigrateCmd is migration cmd command.
func MigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate",
		Short:   "migrates runs specified commands",
		Long:    "migrates runs specified commands",
		PreRunE: loadConfigs,
		RunE:    runMigrateCmd,
		Version: "v0.0.0",
	}

	cmd.Flags().String(configPathFlag, defaultConfigPath, "path to config yml file")

	return cmd
}

const (
	driver = "pgx"
	dir    = "migrations"
)

var errArg = errors.New("command must have at least one argument")

// runMigrateCmd runs the migrate command with arguments.
func runMigrateCmd(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 {
		return errArg
	}

	command := args[0]

	dsn := cfg.Infrastructure.Database.GetDatabaseDSN()

	goose.SetBaseFS(bum.EmbedMigrations)

	db, err := goose.OpenDBWithDriver(driver, dsn)
	if err != nil {
		return fmt.Errorf("goose: failed to open DB: %w", err)
	}

	defer func() {
		if errClose := db.Close(); errClose != nil {
			err = fmt.Errorf("goose: failed to close DB: %w %w", errClose, err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err = goose.RunContext(cmd.Context(), command, db, dir, arguments...); err != nil {
		return fmt.Errorf("failed to run goose command: %w", err)
	}

	return nil
}
