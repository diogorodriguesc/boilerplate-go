package users

import (
	"context"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	usersv1 "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/gen/users/v1"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/handlers"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

var validate = validator.New()

type Server struct {
	usersv1.UnimplementedUsersServiceServer
	api ports.ApiPort
}

func NewServer(api ports.ApiPort) *Server {
	return &Server{api: api}
}

func (s *Server) CreateUser(_ context.Context, req *usersv1.CreateUserRequest) (*usersv1.CreateUserResponse, error) {
	if err := validate.Var(req.GetUsername(), "required,min=3,max=255"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "username: "+err.Error())
	}

	if err := validate.Var(req.GetEmail(), "required,email,max=255"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "email: "+err.Error())
	}

	user, err := s.api.CreateUser(req.GetUsername(), req.GetEmail())
	if err != nil {
		return nil, handlers.MapErrorIntoStatusError(err)
	}

	return &usersv1.CreateUserResponse{User: UserDomainToProto(user)}, nil
}
