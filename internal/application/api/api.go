package api

import (
	"context"
	"strconv"

	"github.com/diogorodriguesc/boilerplate-go/infrastructure/storage"
	usersrepository "github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/users"
	sqlc "github.com/diogorodriguesc/boilerplate-go/internal/adapters/postgres-service-repository/users/tables"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

type Api struct {
	userRepository ports.UserRepository
}

func NewApplication(_ context.Context, pSqlConnection *storage.DB) (ports.ApiPort, func() error, error) {
	return &Api{
		userRepository: usersrepository.New(pSqlConnection, sqlc.New(pSqlConnection.DB)),
	}, func() error { return nil }, nil
}

func (a *Api) SearchUsers(payload ports.SearchUserPayload) ([]domain.UserDomain, error) {
	users, err := a.userRepository.SearchUsers(context.Background(), payload)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (a *Api) CreateUser(username, email string) (*domain.UserDomain, error) {
	user, err := a.userRepository.CreateUser(context.Background(), username, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *Api) GetUserByID(id string) (*domain.UserDomain, error) {
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, applicationerrors.ErrNotFound
	}

	user, err := a.userRepository.GetUserByID(context.Background(), parsedID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
