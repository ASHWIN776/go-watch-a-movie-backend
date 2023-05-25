package main

import (
	"encoding/json"
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
