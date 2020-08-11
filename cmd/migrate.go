package cmd

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// blank because it is needed for migrate, but never used directly
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"leggett.dev/devmarks/api/app"
)

func isValidCommand(command string) bool {
	arr := []string{"up", "down", "drop"}
	for _, a := range arr {
		if a == command {
			return true
		}
	}
	return false
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too many arguments")
		}
		if len(args) != 0 {
			if isValidCommand(args[0]) {
				return nil
			}
			return fmt.Errorf("invalid command specified: %s", args[0])
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetInt("version")

		if version != -1 && len(args) > 0 {
			return errors.New("Cannot use --version flag and a command at the same time")
		}

		a, err := app.New()
		if err != nil {
			return err
		}
		defer a.Close()

		instance, err := postgres.WithInstance(a.Database.DB.DB(), &postgres.Config{})
		if err != nil {
			logrus.Fatal(err)
		}

		m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", instance)
		if err != nil {
			logrus.Fatal(err)
		}

		if version == -1 {
			if args[0] == "up" {
				if err := m.Up(); err != nil {
					logrus.Fatal(err)
				}
				logrus.Info("successfully applied migrations")
			}
			if args[0] == "down" {
				if err := m.Down(); err != nil {
					logrus.Fatal(err)
				}
				logrus.Info("successfully rolled back migrations")
			}
			if args[0] == "drop" {
				if err := m.Drop(); err != nil {
					logrus.Fatal(err)
				}
				logrus.Info("successfully dropped database schema")
			}
		} else {
			if err := m.Migrate(uint(version)); err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("successfully changed to migration version %d", version)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().Int("version", -1, "the migration to run forwards until; if not set, will run all migrations")
}
