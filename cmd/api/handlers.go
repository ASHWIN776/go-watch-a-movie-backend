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

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload

	// validate user against db

	// check password

	// create a jwt user
	u := JWTUser{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
	}

	// Generate Tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
	}

	log.Println(tokens.Token)

	w.Write([]byte(tokens.Token))
}
