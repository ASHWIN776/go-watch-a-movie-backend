package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// Create a router mux
	mux := chi.NewRouter()

	// Use Middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	// Create routes
	mux.Get("/", app.Home)
	mux.Get("/movies", app.AllMovies)

	return mux
}
