package app

import (
	"database/sql"
	"fmt"
	"goBackendServer/internal/api"
	"goBackendServer/internal/middleware"
	"goBackendServer/internal/store"
	"goBackendServer/internal/utils"
	"goBackendServer/migrations"
	"log"
	"net/http"
	"os"
)

// Structure of our Server
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHanlder   *api.TokenHanlder
	Middleware     middleware.UserMiddleware
	DB             *sql.DB
}

// Create a new application
func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	pgDB, err := store.Open()

	if err != nil {
		panic(fmt.Errorf("%s", err))
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	workOutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// handlers
	workOutHandler := api.NewWorkoutHandler(workOutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workOutHandler,
		UserHandler:    userHandler,
		TokenHanlder:   tokenHandler,
		Middleware:     middlewareHandler,
		DB:             pgDB,
	}

	return app, nil
}

// Utility function to health check
func (a Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	a.Logger.Printf("Status is avaiable!\n")
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"status": "OK"})
}
