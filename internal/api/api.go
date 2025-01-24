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

		// check if the user already exists
		if users, err := dbutils.GetUserByUsername(db, user.Username); err != nil || len(users) > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists"})
			return
		}

		// hash the user password
		hashedPsswrd, err := bcrypt.GenerateFromPassword([]byte(
			user.Password), bcrypt.DefaultCost)	
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password"})
		}
		user.Password = string(hashedPsswrd)

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
		var users []models.User
		users, err := dbutils.GetUserByUsername(db, creds.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to query user"})
			return
		} 

		// check if the user exists
		if len(users) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials"})
			return
		} else {
			// check if the password is correct
			if err := bcrypt.CompareHashAndPassword([]byte(users[0].Password),
				[]byte(creds.Password)); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid credentials"})
				return
			}
		}

		// generate a JWT token
		token, err := jwtutils.GenerateJWT(secret, users[0].Username)
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