// internal/sensors/simulateSensor.go - SimulateSensor function
package sensors

import (
	"time"
	"fmt"
	"encoding/json"
	"github.com/kajtekajtek/insight-naturae/pkg/utils" // utils
	"github.com/kajtekajtek/insight-naturae/pkg/models" // models
	mqtt "github.com/eclipse/paho.mqtt.golang" // paho mqtt
	"github.com/google/uuid" // uuid
)

// SimulateSensor simulates a sensor that sends data to a MQTT broker
func SimulateSensor(client mqtt.Client, topics []string, unit string, min, max float64, interval int) {
	// generate a random sensor ID
	id := uuid.New().String()
	for {
		// generate random sensor data
		val := utils.GenerateData(min, max)
		// create a message
		msg := models.SensorData{
			SensorID:  id,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Value:     val,
			Unit:      unit,
		}

		// marshal the message
		payload, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshalling message:", err)
			return
		}

		// publish the message
		for _, t := range topics {
			token := client.Publish(t, 0, false, payload)
			token.Wait()
		}

		// wait
		time.Sleep(time.Duration(interval) * time.Second)
	}

}