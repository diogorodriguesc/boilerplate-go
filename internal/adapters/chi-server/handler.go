package chiserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain/mappers"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
)

func (s *HttpServer) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	email := r.URL.Query().Get("email")
	validate := validator.New()
	if err := validate.Var(email, "required,email"); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, mappers.MapErrorIntoErrorResponse(err))
		return
	}
	user, err := s.api.GetUserByEmail(email)
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
