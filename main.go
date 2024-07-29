package main

import (
	"net/http"
	"time"

	"go-cert-api/initialisations"
	"go-cert-api/routes"
	"go-cert-api/utilitaires"
)

func main() {
	// Load env variables from .env
	initialisations.LoadEnvVariables()

	// Connect to the database
	initialisations.ConnectToDB()

	// Start notification routine
	go utilitaires.PeriodicExpirationChecker(initialisations.DB)
	// Setup Gin router and routes
	router := routes.SetupRouter()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()

}
