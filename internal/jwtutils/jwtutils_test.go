// internal/jwtutils/jwtutils_test.go - jwtutils.go unit tests
package jwtutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/golang-jwt/jwt/v4"
)

const secretString = "testsecret"
const username = "testuser"

func TestGenerateJWT(t *testing.T) {
	secret := []byte(secretString)

	// generate a token
	tokenString, err := GenerateJWT(secret, username)

	assert.NoError(t, err)

	// validate the generated token
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, 
		func(t *jwt.Token) (interface{}, error) { return secret, nil })

	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, username, (*claims)["username"])
}

func TestValidateJWT(t *testing.T) {
	secret := []byte(secretString)

	// generate a token
	tokenString, err := GenerateJWT(secret, username)
	assert.NoError(t, err)

	// validate the token
	validatedUsername, err := ValidateJWT(secret, tokenString)
	assert.NoError(t, err)
	assert.Equal(t, username, validatedUsername)
}

func TestValidateInvalidJWT(t *testing.T) {
	secret := []byte(secretString)
	invalidToken := "invalid.token.string"

	_, err := ValidateJWT(secret, invalidToken)
	assert.Error(t, err)
}

func TestValidateExpiredJWT(t *testing.T) {
	secret := []byte(secretString)

	// generate an expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(-time.Hour).Unix(),
	})
	expiredTokenString, err := token.SignedString(secret)
	assert.NoError(t, err)

	// validate the expired token
	_, err = ValidateJWT(secret, expiredTokenString)
	assert.Error(t, err)
}