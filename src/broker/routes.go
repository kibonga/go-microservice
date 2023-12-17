package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {

	// create new chi router
	mux := chi.NewRouter()

	// determine who can connect
	mux.Use(cors.Handler(cors.Options{
		// accept all origins
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Healthcheck endpoint
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Broker)

	return mux
}
