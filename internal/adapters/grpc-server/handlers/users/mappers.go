package users

import (
	usersv1 "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/gen/users/v1"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
)

func UserDomainToProto(user *domain.UserDomain) *usersv1.User {
	return &usersv1.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
