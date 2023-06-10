package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	Issuer         string
	Audience       string
	Secret         string
	TokenExpiry    time.Duration
	RefreshExpirty time.Duration
	CookieDomain   string
	CookiePath     string
	CookieName     string
}

// Minimal amount of user info needed to issue a token
type JWTUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Tokens are issued as a pair
type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

// Takes a JWTUser instance and generates a TokenPairs instance
func (j *Auth) GenerateTokenPair(user *JWTUser) (TokenPairs, error) {
	// Access Token Gen ---------------------------------
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims (jwt map claims)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// Set the expiry
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Create a signed token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// --------------------------------------------------

	// Refresh Token Gen --------------------------------
	// Create a Refresh Token
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	// Set Claims (jwt map claims)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// Expiry for the Refresh Token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpirty).Unix()

	// Create Signed Refresh Tokens
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// --------------------------------------------------

	// Create TokenPairs instance and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Return TokenPairs instance
	return tokenPairs, nil
}
