package requests

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

type SearchUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=255"`
	Email    string `json:"email" validate:"omitempty,email,max=255"`
}
