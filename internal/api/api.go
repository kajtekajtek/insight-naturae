// internal/api/api.go - API handlers
package api

import (
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/kajtekajtek/insight-naturae/internal/jwtutils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"

	"golang.org/x/crypto/bcrypt"
)

// handle user registration
func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c* gin.Context) {
		var user models.User // user data struct

		// bind the JSON data to the user struct
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// insert the user data into the database
		if err := dbutils.InsertUserData(db, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to register user"})
			return
		}

		// OK
		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully"})
	}
}

// handle user login
func LoginHandler(db *sql.DB, secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds models.User // user credentials struct

		// bind the JSON data to the user credentials struct
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}
		
		// get the user data from the database
		user, err := dbutils.GetUserByUsername(db, creds.Username)
		// check if the user exists and the password is correct
		if err != nil || bcrypt.CompareHashAndPassword(
			[]byte(user.Password), []byte(creds.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials"})
		}

		// generate a JWT token
		token, err := jwtutils.GenerateJWT(secret, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token"})
			return
		}

		// OK
		c.JSON(http.StatusOK, gin.H{
			"token": token})
	}
}