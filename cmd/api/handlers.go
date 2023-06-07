package main

import (
	"log"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go movies api up and running",
		Version: "1.0.0",
	}

	err := app.writeJSON(w, payload, http.StatusOK)
	if err != nil {
		log.Println("Failed to writeJSON err: ", err)
	}
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, movies, http.StatusOK)
	if err != nil {
		log.Println("Failed to writeJSON err: ", err)
	}
}
