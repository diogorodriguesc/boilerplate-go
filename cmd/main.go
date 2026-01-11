package main

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/diogorodriguesc/boilerplate-go/cmd/cli/cobra"
	"github.com/diogorodriguesc/boilerplate-go/cmd/cli/cobra/commands"
	httpserver "github.com/diogorodriguesc/boilerplate-go/cmd/cli/cobra/http-server"
	"github.com/diogorodriguesc/boilerplate-go/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Exit(1)
	}

	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel))

	rootCmd := cobra.GetRootCmd()
	rootCmd.AddCommand(commands.RunDBMigrationsCommand())
	rootCmd.AddCommand(httpserver.ServerHttpCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
