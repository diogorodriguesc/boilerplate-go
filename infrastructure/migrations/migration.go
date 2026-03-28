package migrations

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"

	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
)

const (
	migrationsDir = "migrations"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DBMigrator struct {
	EmbedMigrations fs.FS
	DbConnection    *storage.DB
}

func NewDBMigrator(dbConnection *storage.DB) (*DBMigrator, error) {
	if dbConnection == nil {
		return nil, errors.New("missing internal db")
	}

	return &DBMigrator{
		EmbedMigrations: embedMigrations,
		DbConnection:    dbConnection,
	}, nil
}

func (m *DBMigrator) Migrate(_ context.Context) error {
	goose.SetLogger(CustomLogger{})
	goose.SetBaseFS(m.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose: failed to set goose dialect: %w", err)
	}

	if err := goose.Up(m.DbConnection.DB, migrationsDir); err != nil {
		return fmt.Errorf("goose: failed to up migrations: %w", err)
	}

	return nil
}
