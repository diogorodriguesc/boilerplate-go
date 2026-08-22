package api

import (
	"context"

	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/repository"
	sqlc "github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/tables"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

type Api struct {
	serviceRepository ports.ServiceRepository
}

func NewApplication(_ context.Context, pSqlConnection *storage.DB) (ports.ApiPort, func() error, error) {
	return &Api{
		serviceRepository: repository.NewServiceRepository(pSqlConnection, sqlc.New(pSqlConnection.DB)),
	}, func() error { return nil }, nil
}

func (a *Api) GetUserByEmail(email string) (*domain.UserDomain, error) {
	user, err := a.serviceRepository.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *Api) CreateUser(username, email string) (*domain.UserDomain, error) {
	user, err := a.serviceRepository.CreateUser(context.Background(), username, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
