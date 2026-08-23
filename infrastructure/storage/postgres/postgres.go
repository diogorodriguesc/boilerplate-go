package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/rs/zerolog/log"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
)

func New(ctx context.Context, environment config.Environment, config config.PostgreSQLConfig) (*storage.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.Database,
		),
	)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	schema := "public"
	if environment.IsTesting() {
		schema = "testing"
	}

	log.Ctx(ctx).Info().Str("schema", schema).Msg("Creating schema if not exists")

	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema)
	if err != nil {
		return nil, errors.New("could not create schema")
	}

	log.Ctx(ctx).Info().Str("schema", schema).Msg("Schema created or already exists")

	_, err = db.Exec("SET search_path TO " + schema)
	if err != nil {
		return nil, errors.New("could not set schema")
	}

	log.Ctx(ctx).Info().Str("schema", schema).Msg("Connected to postgres")

	return &storage.DB{DB: db}, nil
}
