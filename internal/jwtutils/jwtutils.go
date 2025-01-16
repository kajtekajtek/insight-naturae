/* internal/jwtutils/jwtutils.go - project specific 
	methods for JWT token management */
package jwtutils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// generate a JWT token with a username and expiration time
func GenerateJWT(secret []byte, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secret)
}

// validate a JWT token and return the username
func ValidateJWT(secret []byte, tokenString string) (string, error) {
	// map to store the claims from the token
	claims := &jwt.MapClaims{}
	// parse the token with the claims and secret
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	// check if the token is valid
	if err != nil || !token.Valid {
		return "", err
	}

	// get the data from the token claims
	username, ok := (*claims)["username"].(string)
	if !ok {
		return "", err
	}

	return username, nil
}