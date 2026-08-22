package ports

import (
	"context"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/go-chi/chi/v5"
)

type (
	SearchUserPayload struct {
		Email string
	}

	ApiPort interface {
		SearchUser(payload SearchUserPayload) ([]domain.UserDomain, error)
		CreateUser(username, email string) (*domain.UserDomain, error)
	}

	HttpService interface {
		Run() error
		Shutdown(ctx context.Context) error
		SetRouter() *chi.Mux
	}

	ServiceRepository interface {
		SearchUsers(ctx context.Context, filters SearchUserPayload) ([]domain.UserDomain, error)
		CreateUser(ctx context.Context, username, email string) (*domain.UserDomain, error)
	}
)
