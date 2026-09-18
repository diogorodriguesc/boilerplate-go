package handlers

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
)

func MapErrorIntoStatusError(err error) error {
	switch {
	case errors.Is(err, applicationerrors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, applicationerrors.ErrDuplicateEntry):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
