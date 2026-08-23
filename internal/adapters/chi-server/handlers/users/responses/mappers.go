package responses

import "github.com/diogorodriguesc/boilerplate-go/internal/application/domain"

func UserDomainToUserResponse(user *domain.UserDomain) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func UserDomainCollectionToUserResponseCollection(users []domain.UserDomain) []UserResponse {
	responseList := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responseList = append(responseList, UserDomainToUserResponse(&user))
	}

	return responseList
}
