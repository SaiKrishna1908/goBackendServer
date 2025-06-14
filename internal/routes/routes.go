package routes

import (
	"goBackendServer/internal/app"

	"github.com/go-chi/chi/v5"
)

// Set up chi routes
func SetUpRoutes(a app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.Authenticate)

		r.Get("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleGetWorkoutById))
		r.Post("/workouts", a.Middleware.RequireUser(a.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleUpdateWorkoutById))
		r.Delete("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleDeleteWorkoutById))
	})

	r.Get("/health", a.HealthCheck)
	r.Post("/users", a.UserHandler.HandleRegisterUser)
	r.Post("/tokens/auth", a.TokenHanlder.HandleCreateToken)

	return r
}
