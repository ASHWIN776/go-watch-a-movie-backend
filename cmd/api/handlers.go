package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
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
	log.Println("Incoming Login Request")

	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against db
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := JWTUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// Generate Tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
	}

	log.Println(tokens.Token)
	refreshCookie := app.auth.GenerateRefreshCookie(tokens.RefreshToken)

	// Will send the cookie with the response
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, tokens, http.StatusAccepted)
}

// Get the refresh token from the cookie sent with the req, get the user id which is stored in the subject of claims and then create another token
func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {

		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.writeJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from token claims
			userId, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.writeJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			// Get the user info from the db using the userID
			user, err := app.DB.GetUserByID(userId)
			if err != nil {
				app.writeJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			// Create Token pairs and sending the access token in json and refresh token in cookie --------------------------------------------------------------------------------------
			// Create a jwt user
			u := JWTUser{
				ID:        userId,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			// Generate Tokens
			tokens, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"))
			}

			refreshCookie := app.auth.GenerateRefreshCookie(tokens.RefreshToken)

			// Will send the cookie with the response
			http.SetCookie(w, refreshCookie)

			// Sending the access token as JSON
			app.writeJSON(w, tokens, http.StatusOK)

			// -----------------------------------------------------------------------------------
		}
	}
}
