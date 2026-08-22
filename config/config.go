package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Environment string

const (
	Development Environment = "dev"
	Testing     Environment = "test"
	Staging     Environment = "stage"
	Production  Environment = "prod"
)

func (e Environment) IsValid() bool {
	switch e {
	case Development, Testing, Staging, Production:
		return true
	}
	return false
}

func (e Environment) IsTesting() bool {
	return e == Testing
}

type (
	Config struct {
		AppName          string      `env:"APP_NAME,required"`
		Env              Environment `env:"ENVIRONMENT,required"`
		LogLevel         int8        `env:"LOG_LEVEL,required"`
		PostgreSQLConfig PostgreSQLConfig
	}

	PostgreSQLConfig struct {
		Host     string `env:"DB_PSQL_HOST,required"`
		Port     string `env:"DB_PSQL_PORT,required"`
		User     string `env:"DB_PSQL_USER,required"`
		Password string `env:"DB_PSQL_PASSWORD,required"`
		Database string `env:"DB_PSQL_DATABASE,required"`
	}
)

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	if !cfg.Env.IsValid() {
		return nil, fmt.Errorf("invalid environment: %s", cfg.Env)
	}

	log.Debug().Interface("env", cfg.Env).Msg("Environment is valid")

	return &cfg, nil
}
