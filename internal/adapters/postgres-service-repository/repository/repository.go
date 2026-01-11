package repository

import (
	"context"

	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/tables"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

type ServiceRepository struct {
	DB      *storage.DB
	queries *sqlc.Queries
}

func NewServiceRepository(db *storage.DB, queries *sqlc.Queries) ports.ServiceRepository {
	return &ServiceRepository{
		DB:      db,
		queries: queries,
	}
}

func (s *ServiceRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserDomain, error) {
	return nil, nil
}
