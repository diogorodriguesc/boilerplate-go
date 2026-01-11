package commands

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/migrations"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage/postgres"
)

func RunDBMigrationsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run-db-migrations",
		Short: "Run db migrations",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("failed to load config: %v", err)
			}

			pSqlConnection, err := postgres.New(cfg.PostgreSQLConfig)
			if err != nil {
				log.Fatalf("failed to connect to postgres: %v", err)
			}

			migrator, err := migrations.NewDBMigrator(pSqlConnection)
			if err := migrator.Migrate(ctx); err != nil {
				log.Fatalf("failed to run migrations: %v", err)
			}
		},
	}
}
