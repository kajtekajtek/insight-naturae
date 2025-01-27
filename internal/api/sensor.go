// internal/api/sensor.go - sensor configuration API handlers
package api

import (
	"database/sql"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
)

// handle adding a sensor for a user to follow
func AddSensorHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userSensor models.UserSensor // user sensor data struct

		// bind the JSON data to the user sensor struct
		if err := c.ShouldBindJSON(&userSensor); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// check if the both fields are provided
		if userSensor.Username == "" || userSensor.SensorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// get user sensors
		if sensors, err := dbutils.GetUserSensors(db, userSensor.Username); 
			err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user sensors"})
			return
		// check if the sensor already exists on the user's list
		} else {
			for _, sensor := range sensors {
				if sensor.SensorID == userSensor.SensorID {
					c.JSON(http.StatusConflict, gin.H{
						"error": "Sensor already exists"})
					return
				}
			}
		}

		// insert the user sensor data into the database
		if err := dbutils.InsertUserSensorData(db, userSensor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add sensor"})
			return
		}

		// OK
		c.JSON(http.StatusCreated, gin.H{
			"message": "Sensor added successfully"})
	}
}

func GetSensorsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Header)
		username := c.GetHeader("Username")
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized"})
			return
		}

		sensors, err := dbutils.GetUserSensors(db, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user sensors"})
			return
		}

		c.JSON(http.StatusOK, sensors)
	}
}

func RemoveSensorHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userSensor models.UserSensor // user sensor data struct

		// bind the JSON data to the user sensor struct
		if err := c.ShouldBindJSON(&userSensor); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// remove the user sensor data from the database
		if err := dbutils.RemoveUserSensor(db, userSensor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to remove sensor"})
			return
		}

		// OK
		c.JSON(http.StatusOK, gin.H{
			"message": "Sensor removed successfully"})
	}
}