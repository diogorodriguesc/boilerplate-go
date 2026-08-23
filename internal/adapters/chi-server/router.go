package chiserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/diogorodriguesc/boilerplate-go/docs"
	usersHandler "github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/handlers/users"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/middlewares"
)

// @title           Golang REST API Boilerplate
// @version         1.0
// @description     This is a sample Go REST API with Swagger documentation.
// @host            localhost:8080
// @BasePath        /
func (s *HttpServer) SetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middlewares.SetJSONResponseMiddleware())

		// Users resources
		r.Post("/users/search", usersHandler.SearchUsers(s.api))
		r.Post("/users", usersHandler.CreateUser(s.api))
	})

	return r
}
