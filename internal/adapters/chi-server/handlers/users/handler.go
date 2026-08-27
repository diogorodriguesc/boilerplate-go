package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/handlers"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/handlers/users/requests"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/handlers/users/responses"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/middlewares"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
)

var validate = validator.New()

// @Summary      Search users
// @Description  Searches for users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      requests.SearchUserRequest  true  "User to search"
// @Success      200    {array}   responses.UserResponse
// @Failure      400    {object}  handlers.ErrorResponse
// @Failure      500    {object}  handlers.ErrorResponse
// @Router       /v1/users/search [post]
func SearchUsers(api ports.ApiPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.SearchUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handlers.MapErrorIntoErrorResponse(err))
			return
		}

		if err := validate.Struct(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handlers.MapErrorIntoErrorResponse(err))
			return
		}

		users, err := api.SearchUsers(ports.SearchUserPayload{
			Username: req.Username,
			Email:    req.Email,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, handlers.MapErrorIntoErrorResponse(err))
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, responses.UserDomainCollectionToUserResponseCollection(users))
	}
}

// @Summary      Create a user
// @Description  Creates a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      requests.CreateUserRequest  true  "User to create"
// @Success      201   {object}  responses.UserResponse
// @Failure      400   {object}  handlers.ErrorResponse
// @Failure      409   {object}  handlers.ErrorResponse
// @Failure      500   {object}  handlers.ErrorResponse
// @Router       /v1/users [post]
func CreateUser(api ports.ApiPort) middlewares.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var req requests.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handlers.MapErrorIntoErrorResponse(err))
			return nil
		}

		if err := validate.Struct(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handlers.MapErrorIntoErrorResponse(err))
			return nil
		}

		user, err := api.CreateUser(req.Username, req.Email)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, responses.UserDomainToUserResponse(user))
		return nil
	}
}

// @Summary      Get a user by ID
// @Description  Retrieves a user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  responses.UserResponse
// @Failure      404  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /v1/users/{id} [get]
func GetUser(api ports.ApiPort) middlewares.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := chi.URLParam(r, "id")
		user, err := api.GetUserByID(id)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, responses.UserDomainToUserResponse(user))
		return nil
	}
}

// @Summary      Delete a user by ID
// @Description  Deletes a user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path  string  true  "User ID"
// @Success      204
// @Failure      404  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /v1/users/{id} [delete]
func DeleteUser(api ports.ApiPort) middlewares.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := chi.URLParam(r, "id")
		if err := api.DeleteUser(id); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

// @Summary      List users
// @Description  Retrieves a paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "Page number"       default(1)
// @Param        pageSize  query     int  false  "Items per page"    default(20)
// @Success      200  {object}  responses.ListUsersResponse
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /v1/users [get]
func ListUsers(api ports.ApiPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := requests.ListUsersRequest{
			Page:     defaultPage,
			PageSize: defaultPageSize,
		}

		if page := r.URL.Query().Get("page"); page != "" {
			parsedPage, err := strconv.Atoi(page)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, handlers.MapErrorIntoErrorResponse(errors.New("page must be an integer")))
				return
			}
			req.Page = parsedPage
		}

		if pageSize := r.URL.Query().Get("pageSize"); pageSize != "" {
			parsedPageSize, err := strconv.Atoi(pageSize)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, handlers.MapErrorIntoErrorResponse(errors.New("pageSize must be an integer")))
				return
			}
			req.PageSize = parsedPageSize
		}

		if err := validate.Struct(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, handlers.MapErrorIntoErrorResponse(err))
			return
		}

		users, total, err := api.ListUsers(req.Page, req.PageSize)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, handlers.MapErrorIntoErrorResponse(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, responses.UserDomainCollectionToListUsersResponse(users, req.Page, req.PageSize, total))
	}
}
