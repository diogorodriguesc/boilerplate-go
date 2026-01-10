package mappers

import (
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain/responses"
)

func MapUserDomainIntoUserResponse(user domain.UserDomain) responses.UserResponse {
	return responses.UserResponse{
		Email: user.Email,
	}
}

func MapErrorIntoErrorResponse(err error) responses.ErrorResponse {
	return responses.ErrorResponse{
		Error: responses.ErrorDetail{
			Message: err.Error(),
		},
	}
}
