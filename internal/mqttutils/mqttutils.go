// internal/mqttutils/mqttutils.go - project specific methods for handling MQTT communication
package mqttutils

import (
	"log"
	"encoding/json"
	"database/sql"

	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/internal/websocket"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/* function MessageHandler is a closure that returns a function 
	which inserts the received sensor data into the database and
	sends it through WebSocket to sensor's connected subscribers */
func MessageHandler(db *sql.DB, cm *websocket.WSClientManager) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		// parse the message payload into a SensorData struct
		data := models.SensorData{}
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Failed to parse the message payload: %v", err)
			return
		}

		// insert the data into the database
		if err := dbutils.InsertSensorData(db, data); err != nil {
			log.Printf("Failed to insert sensor data into the database: %v", err)
			return
		}

		// send the data through WebSocket to sensor's subscribers
		if subscribers, exists := cm.Subscriptions[data.SensorID]; exists {
			message, _ := json.Marshal(data)
			for username := range subscribers {
				err := cm.SendMessage(username, message)
				if err != nil {
					log.Printf("Failed to send sensor data to client %s: %v",
						username, err)
				}
			}
		}
	}
}