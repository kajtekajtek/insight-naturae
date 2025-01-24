// pkg/models/UserSensor.go - data model for sensor followed by user
package models

type UserSensor struct {
	Username string `json:"username"`
	SensorID string `json:"sensor_id"`
}