/* pkg/models/SensorData.go - defines the data structure 
	for the sensor readings */
package models

type SensorData struct {
	SensorID  string  `json:"sensor_id"`
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}