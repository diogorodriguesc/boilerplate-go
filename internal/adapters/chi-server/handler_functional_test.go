//go:build functional

package chiserver_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/migrations"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage/postgres"
	chiserver "github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/api"
)

const (
	postgresDockerImage = "postgres:16-alpine"

	postgresFixtureReset    = "postgres/000_reset.sql"
	postgresFixtureAddUsers = "postgres/100_add_users.sql"
)

var (
	postgresDatabaseConnection *sql.DB
	postgresStorage            *storage.DB

	postgresContainer testcontainers.Container
)

func TestMain(m *testing.M) {
	os.Exit(runFunctionalTestMain(m))
}

func performRequest(t *testing.T, email string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/users?email=%s", email), nil)
	req.Header.Set("Content-Type", "application/json")

	application, _, err := api.NewApplication(t.Context(), postgresStorage)
	if err != nil {
		t.Fatal(err)
	}

	httpServer := chiserver.NewHttpServer(t.Context(), application)

	recorder := httptest.NewRecorder()
	httpServer.SetRouter().ServeHTTP(recorder, req)
	return recorder
}

func runFunctionalTestMain(m *testing.M) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var err error
	var migrator *migrations.DBMigrator
	postgresContainer, migrator, err = startPostgresContainer(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres container: %v\n", err)
		return 1
	}

	if err = migrator.Migrate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run migrator: %v\n", err)
		return 1
	}

	code := m.Run()
	terminateContainers()
	return code
}

func startPostgresContainer(ctx context.Context) (testcontainers.Container, *migrations.DBMigrator, error) {
	req := testcontainers.ContainerRequest{
		Image:        postgresDockerImage,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres", // ggignore
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
		return nil, nil, err
	}

	db, host, port, err := getDatabaseConnection(ctx, container)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}

	if err := waitForPing(ctx, db); err != nil {
		_ = db.Close()
		_ = container.Terminate(ctx)
		return nil, nil, err
	}

	postgresStorage, err = postgres.New(
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

	migrator, err := migrations.NewDBMigrator(postgresStorage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create migrator: %v\n", err)
		return nil, nil, err
	}

	return container, migrator, nil
}

func getDatabaseConnection(ctx context.Context, container testcontainers.Container) (*sql.DB, string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", "", err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", "", err
	}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host,
			port.Port(),
			"postgres",
			"postgres",
			"postgres",
		),
	)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", "", err
	}
	return db, host, port.Port(), nil
}

func seedCase(t *testing.T, postgresFixture string) error {
	t.Helper()
	db, _, _, err := getDatabaseConnection(t.Context(), postgresContainer)
	if err != nil {
		return nil
	}
	require.NoError(t, applyFixture(t, db, postgresFixtureReset))
	require.NoError(t, applyFixture(t, db, postgresFixture))
	return nil
}

func fixtureAbsPath(relative string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
	return filepath.Join(projectRoot, "internal", "adapters", "chi-server", "testdata", relative)
}

func applyFixture(t *testing.T, db *sql.DB, fixtureRelativePath string) error {
	t.Helper()
	content, err := os.ReadFile(fixtureAbsPath(fixtureRelativePath))
	if err != nil {
		return err
	}
	_, err = db.Exec(string(content))
	return err
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

func terminateContainers() {
	if postgresDatabaseConnection != nil {
		_ = postgresDatabaseConnection.Close()
	}

	if postgresContainer != nil {
		_ = postgresContainer.Terminate(context.Background())
	}
}

func TestHandler_Functional_GetUserByEmail(t *testing.T) {
	require.NoError(t, seedCase(t, postgresFixtureAddUsers))
	recorder := performRequest(t, "foo@gmail.com")
	require.Equal(t, recorder.Code, http.StatusFound)
	require.JSONEq(t, `{
	"id": 1,
	"username": "foo",
	"email": "foo@gmail.com"
}`, recorder.Body.String())
}
