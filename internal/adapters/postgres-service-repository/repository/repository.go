package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/tables"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

const pqUniqueViolationCode = "23505"

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

func (s *ServiceRepository) SearchUsers(ctx context.Context, filters ports.SearchUserPayload) ([]domain.UserDomain, error) {
	email := sql.NullString{String: filters.Email, Valid: filters.Email != ""}

	users, err := s.queries.SearchUsers(ctx, email)
	if err != nil {
		return nil, err
	}

	results := make([]domain.UserDomain, 0, len(users))
	for _, user := range users {
		results = append(results, domain.UserDomain{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return results, nil
}

func (s *ServiceRepository) CreateUser(ctx context.Context, username, email string) (*domain.UserDomain, error) {
	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: username,
		Email:    email,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqUniqueViolationCode {
			return nil, applicationerrors.ErrDuplicateEntry
		}
		return nil, err
	}

	return &domain.UserDomain{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
