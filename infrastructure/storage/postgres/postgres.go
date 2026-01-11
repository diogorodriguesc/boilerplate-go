package postgres

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"time"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
)

func New(config config.PostgreSQLConfig) (*storage.DB, error) {
	db, err := sql.Open(
		"nrpostgres",
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
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &storage.DB{DB: db}, nil
}
