package app

import (
	"fmt"
	"goBackendServer/internal/api"
	"log"
	"net/http"
	"os"
)

// Structure of our Server
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

// Create a new application
func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// TODO: data stores

	// handlers
	workOutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workOutHandler,
	}

	return app, nil
}

// Utility function to health check
func (a Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available!\n")
}
