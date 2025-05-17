package routes

import (
	"goBackendServer/internal/app"

	"github.com/go-chi/chi/v5"
)

// Set up chi routes
func SetUpRoutes(a app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", a.HealthCheck)
	r.Get("/workouts/{id}", a.WorkoutHandler.HandleWorkoutHandler)
	r.Post("/workouts", a.WorkoutHandler.HandleCreateWorkout)

	return r
}
