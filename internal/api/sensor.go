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

// handle sensor subscription requests
func SubscribeSensorHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sensorSubscription models.SensorSubscription // user sensor data struct

		// bind the JSON data to the user sensor struct
		if err := c.ShouldBindJSON(&sensorSubscription); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// check if the both fields are provided
		if sensorSubscription.Username == "" || sensorSubscription.SensorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// get user sensors
		if sensors, err := dbutils.GetUserSubscriptions(db, sensorSubscription.Username); 
			err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user sensors"})
			return
		// check if the sensor already exists on the user's list
		} else {
			for _, sensor := range sensors {
				if sensor.SensorID == sensorSubscription.SensorID {
					c.JSON(http.StatusConflict, gin.H{
						"error": "Sensor already exists"})
					return
				}
			}
		}

		// insert the user sensor data into the database
		if err := dbutils.InsertSensorSubscriptionData(db, sensorSubscription); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add sensor"})
			return
		}

		// OK
		c.JSON(http.StatusCreated, gin.H{
			"message": "Sensor added successfully"})
	}
}

// get user subscribed sensors
func GetSensorsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Header)
		username := c.GetHeader("Username")
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized"})
			return
		}

		sensors, err := dbutils.GetUserSubscriptions(db, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user sensors"})
			return
		}

		c.JSON(http.StatusOK, sensors)
	}
}

// handle sensor unsubscription requests
func UnsubscribeSensorHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sensorSubscription models.SensorSubscription // user sensor data struct

		// bind the JSON data to the user sensor struct
		if err := c.ShouldBindJSON(&sensorSubscription); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// check if the both fields are provided
		if sensorSubscription.Username == "" || sensorSubscription.SensorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// check if the sensor exists on the user's list
		if sensors, err := dbutils.GetUserSubscriptions(db, sensorSubscription.Username);
			err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get user sensors"})
			return
		} else {
			found := false
			for _, sensor := range sensors {
				if sensor.SensorID == sensorSubscription.SensorID {
					found = true
					break
				}
			}
			if !found {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Sensor not found"})
				return
			}
		}

		// remove the user's sensor subscription from the database
		if err := dbutils.RemoveSensorSubscription(db, sensorSubscription); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to remove sensor"})
			return
		}

		// OK
		c.JSON(http.StatusOK, gin.H{
			"message": "Sensor removed successfully"})
	}
}