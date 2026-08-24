package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	sqlc "github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/users/tables"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

const pqUniqueViolationCode = "23505"

type Repository struct {
	DB      *storage.DB
	queries *sqlc.Queries
}

func New(db *storage.DB, queries *sqlc.Queries) ports.UserRepository {
	return &Repository{
		DB:      db,
		queries: queries,
	}
}

func (s *Repository) SearchUsers(ctx context.Context, filters ports.SearchUserPayload) ([]domain.UserDomain, error) {
	users, err := s.queries.SearchUsers(ctx, sqlc.SearchUsersParams{
		Email:    sql.NullString{String: filters.Email, Valid: filters.Email != ""},
		Username: sql.NullString{String: filters.Username, Valid: filters.Username != ""},
	})
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

func (s *Repository) CreateUser(ctx context.Context, username, email string) (*domain.UserDomain, error) {
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

func (s *Repository) GetUserByID(ctx context.Context, id int64) (*domain.UserDomain, error) {
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, applicationerrors.ErrNotFound
		}
		return nil, err
	}

	return &domain.UserDomain{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
