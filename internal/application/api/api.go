package api

import (
	"context"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

type Api struct{}

func NewApplication(_ context.Context) (ports.ApiPort, func() error, error) {
	return &Api{}, func() error { return nil }, nil
}

func (a *Api) GetUserByEmail(email string) (domain.UserDomain, error) {
	return domain.UserDomain{
		Email: email,
	}, nil
}
