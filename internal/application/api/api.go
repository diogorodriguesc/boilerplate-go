package api

import (
	"context"
	"os"

	"github.com/diogorodriguesc/boilerplate-go/config"
	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage/postgres"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/repository"
	sqlc "github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/tables"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

type Api struct {
	serviceRepository ports.ServiceRepository
}

func NewApplication(_ context.Context) (ports.ApiPort, func() error, error) {
	cfg, err := config.Load()
	if err != nil {
		os.Exit(1)
	}

	pSqlConnection, err := postgres.New(cfg.PostgreSQLConfig)
	if err != nil {
		return nil, nil, err
	}

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
