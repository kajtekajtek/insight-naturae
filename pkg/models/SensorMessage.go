/* pkg/models/SensorMessage.go - defines the data structure for the messages
	sent by the sensors to the MQTT broker */
package models

type SensorMessage struct {
	SensorID  string  `json:"sensor_id"`
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}