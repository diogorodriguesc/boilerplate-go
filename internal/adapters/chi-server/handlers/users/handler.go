package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain/mappers"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

// @Summary      Get user by email
// @Description  Retrieves a user's details by their email address
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email  query     string  true  "User email address"
// @Success      302    {object}  responses.UserResponse
// @Failure      400    {object}  responses.ErrorResponse
// @Failure      404    {object}  responses.ErrorResponse
// @Failure      500    {object}  responses.ErrorResponse
// @Router       /v1/users [get]
func GetUserByEmail(api ports.ApiPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		validate := validator.New()
		if err := validate.Var(email, "required,email"); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, mappers.MapErrorIntoErrorResponse(err))
			return
		}
		user, err := api.GetUserByEmail(email)
		if err != nil {
			if errors.Is(err, applicationerrors.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			render.JSON(w, r, mappers.MapErrorIntoErrorResponse(err))
			return
		}
		w.WriteHeader(http.StatusFound)
		render.JSON(w, r, mappers.MapUserDomainIntoUserResponse(user))
	}
}
