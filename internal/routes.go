package routes

import (
	"goBackendServer/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetUpRoutes(a app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", a.HealthCheck)

	return r
}
