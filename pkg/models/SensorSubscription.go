// pkg/models/SensorSubscription.go - data model for sensor followed by user
package models

type SensorSubscription struct {
	Username string `json:"username"`
	SensorID string `json:"sensor_id"`
}