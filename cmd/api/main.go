package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 8000

type application struct {
	Domain string
}

func main() {
	// set application config
	var app application

	// Read command line flags

	// Connect to db

	app.Domain = "example.com"

	// Start a web server
	log.Printf("Listening on %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
