package grpcserver

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage/postgres"
	grpcServerAdapter "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/api"
)

const (
	CommandUse   = "grpc-server"
	CommandShort = "Run gRPC Server"
)

func ServerGrpcCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandUse,
		Short: CommandShort,
		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()
			ctx = log.Logger.WithContext(ctx)

			cfg, err := config.Load()
			if err != nil {
				os.Exit(1)
			}

			pSqlConnection, err := postgres.New(ctx, cfg.Env, cfg.PostgreSQLConfig)
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
				grpcServer:  grpcServerAdapter.NewGrpcServer(ctx, application),
			}
			container.start(ctx)
		},
	}
}
