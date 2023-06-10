package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8000

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func main() {
	// set application config
	var app application

	// Read command line flags
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=54320 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.Parse()

	// Connect to db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// Will add the connection to the PostgresDBRepo struct, and then assign that whole thing to app.DB (will satisfy the repository.DatabaseRepo interface)
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	// Defer - Close connection
	defer app.DB.Connection().Close()

	app.auth = Auth{
		Issuer:         app.JWTIssuer,
		Audience:       app.JWTAudience,
		Secret:         app.JWTSecret,
		TokenExpiry:    time.Minute * 15,
		RefreshExpirty: time.Hour * 24,
		CookiePath:     "/",
		CookieName:     "__Host-refresh_token",
		CookieDomain:   app.CookieDomain,
	}

	// Start a web server
	log.Printf("Listening on %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
