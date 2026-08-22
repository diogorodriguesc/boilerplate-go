// Package tests provides shared helpers for functional tests that need a
// real Postgres instance, so individual handler test packages don't have
// to duplicate testcontainers setup/teardown logic.
package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/migrations"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage/postgres"
)

const postgresDockerImage = "postgres:16-alpine"

// PostgresEnvironment is a running Postgres testcontainer together with the
// application storage connected to it.
type PostgresEnvironment struct {
	Storage *storage.DB

	container testcontainers.Container
	rawDB     *sql.DB
}

// StartPostgresEnvironment starts a Postgres container, connects the
// application storage to it and runs all pending migrations.
func StartPostgresEnvironment(ctx context.Context) (*PostgresEnvironment, error) {
	req := testcontainers.ContainerRequest{
		Image:        postgresDockerImage,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(1).
			WithStartupTimeout(2 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	env := &PostgresEnvironment{container: container}

	host, port, err := containerAddress(ctx, container)
	if err != nil {
		env.Close()
		return nil, err
	}

	env.rawDB, err = openDB(host, port)
	if err != nil {
		env.Close()
		return nil, err
	}

	if err := waitForPing(ctx, env.rawDB); err != nil {
		env.Close()
		return nil, err
	}

	env.Storage, err = postgres.New(
		ctx,
		config.Testing,
		config.PostgreSQLConfig{
			Host:     host,
			Port:     port,
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
		},
	)
	if err != nil {
		env.Close()
		return nil, err
	}

	migrator, err := migrations.NewDBMigrator(env.Storage)
	if err != nil {
		env.Close()
		return nil, err
	}

	if err := migrator.Migrate(ctx); err != nil {
		env.Close()
		return nil, err
	}

	return env, nil
}

// Close releases the raw connection and terminates the container.
func (e *PostgresEnvironment) Close() {
	if e.rawDB != nil {
		_ = e.rawDB.Close()
	}
	if e.container != nil {
		_ = e.container.Terminate(context.Background())
	}
}

// ApplyFixture executes the SQL statements contained in the file at path
// against the environment's database.
func (e *PostgresEnvironment) ApplyFixture(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = e.rawDB.Exec(string(content))
	return err
}

// RunFunctionalMain starts a Postgres environment, hands it to onReady so
// the caller's TestMain can expose it to its tests, runs m.Run() and tears
// the environment down afterwards.
func RunFunctionalMain(m *testing.M, onReady func(*PostgresEnvironment)) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	env, err := StartPostgresEnvironment(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres environment: %v\n", err)
		return 1
	}
	defer env.Close()

	onReady(env)

	return m.Run()
}

func containerAddress(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return "", "", err
	}

	return host, port.Port(), nil
}

func openDB(host, port string) (*sql.DB, error) {
	return sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, "postgres", "postgres", "postgres",
		),
	)
}

func waitForPing(ctx context.Context, db *sql.DB) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	deadline := time.After(60 * time.Second)
	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline:
			return errors.New("timed out waiting for ping")
		case <-ticker.C:
		}
	}
}
