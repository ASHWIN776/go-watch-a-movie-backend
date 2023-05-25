package main

import (
	"backend/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
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

	// Convert the payload into json
	out, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return
	}

	// Specify Headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	rd, _ := time.Parse("2006-01-02", "1986-03-07")

	movie := models.Movie{
		ID:          1,
		Title:       "Highlander",
		ReleaseDate: rd,
		Runtime:     116,
		MpaaRating:  "R",
		Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Iste enim veritatis aliquam iusto consequatur amet nisi dolorum delectus esse error deleniti eligendi rem debitis, vel quod. Incidunt quae dolorum modi.",
	}

	movies = append(movies, movie)

	rd, _ = time.Parse("2006-01-02", "2012-03-07")

	movie = models.Movie{
		ID:          2,
		Title:       "Avenger",
		ReleaseDate: rd,
		Runtime:     116,
		MpaaRating:  "R",
		Description: "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Iste enim veritatis aliquam iusto consequatur amet nisi dolorum delectus esse error deleniti eligendi rem debitis, vel quod. Incidunt quae dolorum modi.",
	}

	movies = append(movies, movie)

	// Convert the payload to json
	out, err := json.Marshal(movies)
	if err != nil {
		log.Println(err)
		return
	}

	// Specify Headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

}
