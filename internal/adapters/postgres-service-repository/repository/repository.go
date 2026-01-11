package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/tables"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
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
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, applicationerrors.ErrRecordNotFound
		}
		return nil, err
	}

	return &domain.UserDomain{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
