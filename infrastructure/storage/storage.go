package storage

import (
	"database/sql"
)

type DB struct {
	DB *sql.DB
}
