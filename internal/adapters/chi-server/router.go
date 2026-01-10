package chiserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/middlewares"
)

func (s *HttpServer) setRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middlewares.SetJSONResponseMiddleware())

		r.Get("/v1/users", s.GetUserByEmail)
	})

	return r
}
