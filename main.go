package main

import (
	"flag"
	"fmt"
	"goBackendServer/internal/app"
	"goBackendServer/internal/routes"
	"net/http"
	"time"
)

/*
starts the go server
*/
func main() {
	var port int

	// take input from command line if not present fall back to 8080
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	// Initialize new application
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	// set up routes
	r := routes.SetUpRoutes(*app)

	// server properties
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}

	app.Logger.Printf("Started go server at %d", port)

	// Start the server
	err = server.ListenAndServe()

	if err != nil {
		app.Logger.Fatal(err)
	}
}
