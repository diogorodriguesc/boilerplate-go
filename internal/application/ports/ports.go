package ports

import (
	"context"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/go-chi/chi/v5"
)

type (
	SearchUserPayload struct {
		Username string
		Email    string
	}

	ApiPort interface {
		SearchUsers(payload SearchUserPayload) ([]domain.UserDomain, error)
		CreateUser(username, email string) (*domain.UserDomain, error)
		GetUserByID(id string) (*domain.UserDomain, error)
		ListUsers(page, pageSize int) ([]domain.UserDomain, int64, error)
		DeleteUser(id string) error
	}

	HttpService interface {
		Run() error
		Shutdown(ctx context.Context) error
		SetRouter() *chi.Mux
	}

	UserRepository interface {
		SearchUsers(ctx context.Context, filters SearchUserPayload) ([]domain.UserDomain, error)
		CreateUser(ctx context.Context, username, email string) (*domain.UserDomain, error)
		GetUserByID(ctx context.Context, id int64) (*domain.UserDomain, error)
		ListUsers(ctx context.Context, limit, offset int32) ([]domain.UserDomain, error)
		CountUsers(ctx context.Context) (int64, error)
		DeleteUser(ctx context.Context, id int64) error
	}
)
