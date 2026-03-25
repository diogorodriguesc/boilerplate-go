package httpserver

import (
	"context"
	"os"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage/postgres"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	chiserver "github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/api"
)

const (
	CommandUse   = "http-server"
	CommandShort = "Run Http Server"
)

func ServerHttpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandUse,
		Short: CommandShort,
		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()

			cfg, err := config.Load()
			if err != nil {
				os.Exit(1)
			}

			pSqlConnection, err := postgres.New(cfg.PostgreSQLConfig)
			if err != nil {
				os.Exit(1)
			}

			application, closeDbConnection, err := api.NewApplication(ctx, pSqlConnection)
			if err != nil {
				log.Error().Err(err).Msg("failed to initialize application")
			}
			defer func() {
				if err := closeDbConnection(); err != nil {
					log.Error().Err(err).Msg("failed to close database connections")
				}
			}()
			container := Container{
				application: application,
				httpServer:  chiserver.NewHttpServer(ctx, application),
			}
			container.start(ctx)
		},
	}
}
