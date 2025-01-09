// internal/mqttutils/mqttutils.go - project specific methods for handling MQTT communication
package mqttutils

import (
	"fmt"
	"encoding/json"
	"database/sql"
	"strings"

	"github.com/kajtekajtek/insight-naturae/pkg/utils"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

/* loads the MQTT connection options from the environment 
	variables and returns them as a MqttConn struct */
func LoadConnOpts() models.MqttConn {
	var c models.MqttConn

	// load the environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file; continuing with the default values...")
	}

	c.Scheme = utils.Getenv("MQTT_SCHEME", "tcp")
	c.Host = utils.Getenv("MQTT_HOST", "localhost")
	c.Port = utils.Getenv("MQTT_PORT", "1883")
	topics := utils.Getenv("MQTT_TOPICS", "insight-naturae/sensors")
	c.Topics = strings.Split(topics, ":")

	return c
}

/* function MessageHandler is a closure that returns a function 
	which inserts the received sensor data into the database */
func MessageHandler(db *sql.DB) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		// parse the message payload into a SensorData struct
		data := models.SensorData{}
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}

		// insert the data into the database
		if err := dbutils.InsertSensorData(db, data); err != nil {
			fmt.Println("Failed to insert sensor data into the database:", err)
			return
		}
	}
}