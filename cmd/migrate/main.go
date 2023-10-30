package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"log/slog"

	"github.com/fiber-bot/cmd/migrate/migrations"
	"github.com/fiber-bot/config"
	"github.com/uptrace/bun/migrate"

	"github.com/urfave/cli/v2"
)

func main() {
	envConfig, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalln(err)
	}

	logger := config.InitLogger(envConfig)

	db, err := config.InitDB(envConfig)
	if err != nil {
		logger.Error("Database connection error", slog.Any("error", err))
		panic("database error")
	}

	app := &cli.App{
		Name:                 "Fiber-Bot",
		EnableBashCompletion: true, Suggest: true,
		Commands: []*cli.Command{
			dbCommand(migrate.NewMigrator(db, migrations.Migrations), *logger),
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.Error("Cli error", slog.Any("error", err))
		os.Exit(1)
	}
}

func dbCommand(migrator *migrate.Migrator, logger slog.Logger) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context)

					group, err := migrator.Migrate(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						logger.Info("there are no new migrations to run (database is up to date)")
						return nil
					}
					logger.Info(fmt.Sprintf("migrated to %s", group))
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context)

					group, err := migrator.Rollback(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						logger.Info("there are no groups to roll back")
						return nil
					}
					logger.Info(fmt.Sprintf("rolled back %s", group))
					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(c.Context, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						logger.Info(fmt.Sprintf("created migration %s (%s)", mf.Name, mf.Path))
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ms, err := migrator.MigrationsWithStatus(c.Context)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())
					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					group, err := migrator.Migrate(c.Context, migrate.WithNopMigration())
					if err != nil {
						return err
					}
					if group.IsZero() {
						logger.Info("there are no new migrations to mark as applied")
						return nil
					}
					fmt.Printf("marked as applied %s", group)
					return nil
				},
			},
		},
	}
}
