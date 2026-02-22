package http

import (
	"Assignment2/internal/handlers"
	"Assignment2/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(h *handlers.Handler, apiKey string, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger())

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Group(func(pr chi.Router) {
		pr.Use(middleware.APIKey(apiKey))

		pr.Get("/healthz", h.Health)
		pr.Post("/auth/login", h.Login())

		pr.Route("/users", func(ur chi.Router) {
			ur.Use(middleware.JWTAuth(jwtSecret))

			ur.Get("/", h.GetUsers)
			ur.Get("/{id}", h.GetUserByID)
			ur.Post("/", h.CreateUser)
			ur.Put("/{id}", h.UpdateUser)
			ur.Patch("/{id}", h.UpdateUser)
			ur.Delete("/{id}", h.DeleteUser)
		})
	})

	return r
}
