package routes

import (
	"goBackendServer/internal/app"

	"github.com/go-chi/chi/v5"
)

// Set up chi routes
func SetUpRoutes(a app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", a.HealthCheck)

	r.Get("/workouts/{id}", a.WorkoutHandler.HandleGetWorkoutById)
	r.Post("/workouts", a.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", a.WorkoutHandler.HandleUpdateWorkoutById)
	r.Delete("/workouts/{id}", a.WorkoutHandler.HandleDeleteWorkoutById)

	r.Post("/users", a.UserHandler.HandleRegisterUser)

	return r
}
