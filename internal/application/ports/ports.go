package ports

import (
	"context"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
)

type (
	ApiPort interface {
		GetUserByEmail(email string) (*domain.UserDomain, error)
	}

	HttpService interface {
		Run() error
		Shutdown(ctx context.Context) error
	}

	ServiceRepository interface {
		GetUserByEmail(ctx context.Context, email string) (*domain.UserDomain, error)
	}
)
