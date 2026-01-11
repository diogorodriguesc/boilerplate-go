package chiserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/diogorodriguesc/boilerplate-go/internal/application/domain/mappers"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
	"github.com/go-chi/render"
)

func (s *HttpServer) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := s.api.GetUserByEmail(r.URL.Query().Get("email"))
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
