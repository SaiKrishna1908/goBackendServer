package app

import (
	"database/sql"
	"fmt"
	"goBackendServer/internal/api"
	"goBackendServer/internal/store"
	"log"
	"net/http"
	"os"
)

// Structure of our Server
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

// Create a new application
func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pgDB, err := store.Open()

	if err != nil {
		panic(fmt.Errorf("%s", err))
	}

	// TODO: data stores

	// handlers
	workOutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workOutHandler,
		DB:             pgDB,
	}

	return app, nil
}

// Utility function to health check
func (a Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available!\n")
}
