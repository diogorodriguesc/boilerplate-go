package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/handlers/users/requests"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain/mappers"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

var validate = validator.New()

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

// @Summary      Create a user
// @Description  Creates a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      requests.CreateUserRequest  true  "User to create"
// @Success      201   {object}  responses.UserResponse
// @Failure      400   {object}  responses.ErrorResponse
// @Failure      409   {object}  responses.ErrorResponse
// @Failure      500   {object}  responses.ErrorResponse
// @Router       /v1/users [post]
func CreateUser(api ports.ApiPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, mappers.MapErrorIntoErrorResponse(err))
			return
		}

		if err := validate.Struct(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, mappers.MapErrorIntoErrorResponse(err))
			return
		}

		user, err := api.CreateUser(req.Username, req.Email)
		if err != nil {
			if errors.Is(err, applicationerrors.ErrDuplicateEntry) {
				w.WriteHeader(http.StatusConflict)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			render.JSON(w, r, mappers.MapErrorIntoErrorResponse(err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, mappers.MapUserDomainIntoUserResponse(user))
	}
}
