// internal/api/sensor.go - sensor configuration API handlers
package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/internal/wsutils"
)

/* Request body should contain only the sensor ID as the username is
	taken from the token */
type subscribeSensorRequestBody struct {
	SensorID string `json:"sensor_id"`
}

// handle sensor subscription requests
func SubscribeSensorHandler(db *sql.DB, cm *wsutils.WSClientManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// sensor subscription data struct
		var sensorSubscription models.SensorSubscription
		// expected request body
		var requestBody subscribeSensorRequestBody

		// bind the JSON data to the expected request body
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}
		
		// get the username retrieved from the token
		username := c.GetHeader("Username")

		// check if both fields are provided
		if username == "" || requestBody.SensorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// set the sensor subscription data
		sensorSubscription.Username = username
		sensorSubscription.SensorID = requestBody.SensorID

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
						"error": "Sensor already subscribed"})
					return
				}
			}
		}

		// insert the user sensor data into the database
		if err := dbutils.InsertSensorSubscriptionData(db, sensorSubscription); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to subscribe sensor"})
			return
		}

		// add subscription to the WebSocket client manager
		cm.Subscribe(sensorSubscription.Username, sensorSubscription.SensorID)

		// OK
		c.JSON(http.StatusCreated, gin.H{
			"message": "Sensor subscribed successfully"})
	}
}

// get user subscribed sensors
func GetSensorsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		// sensor subscription data struct
		var sensorSubscription models.SensorSubscription
		// expected request body
		var requestBody subscribeSensorRequestBody

		// bind the JSON data to the user sensor struct
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// get the username retrieved from the token
		username := c.GetHeader("Username")

		// check if both fields are provided
		if username == "" || requestBody.SensorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input"})
			return
		}

		// set the sensor subscription data
		sensorSubscription.Username = username
		sensorSubscription.SensorID = requestBody.SensorID

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
				"error": "Failed to unsubscribe sensor"})
			return
		}

		// OK
		c.JSON(http.StatusOK, gin.H{
			"message": "Sensor unsubscribed successfully"})
	}
}