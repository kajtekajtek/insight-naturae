// internal/mqttutils/mqttutils.go - project specific methods for handling MQTT communication
package mqttutils

import (
	"fmt"
	"encoding/json"
	"database/sql"

	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/* function MessageHandler is a closure that returns a function 
	which inserts the received sensor data into the database */
func MessageHandler(db *sql.DB) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		// parse the message payload into a SensorData struct
		data := models.SensorData{}
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			fmt.Println("Failed while unmarshalling message:", err)
			return
		}

		// insert the data into the database
		if err := dbutils.InsertSensorData(db, data); err != nil {
			fmt.Println("Failed to insert sensor data into the database:", err)
			return
		}
	}
}