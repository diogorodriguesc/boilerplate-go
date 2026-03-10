package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/rs/zerolog/log"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
)

func New(config config.PostgreSQLConfig) (*storage.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
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

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Info().Msg("Connected to postgres")

	return &storage.DB{DB: db}, nil
}
