package mappers

import (
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain/responses"
)

func MapUserDomainIntoUserResponse(user *domain.UserDomain) responses.UserResponse {
	return responses.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func MapUserDomainsIntoUserResponses(users []domain.UserDomain) []responses.UserResponse {
	responseList := make([]responses.UserResponse, 0, len(users))
	for _, user := range users {
		responseList = append(responseList, MapUserDomainIntoUserResponse(&user))
	}

	return responseList
}

func MapErrorIntoErrorResponse(err error) responses.ErrorResponse {
	return responses.ErrorResponse{
		Error: responses.ErrorDetail{
			Message: err.Error(),
		},
	}
}
