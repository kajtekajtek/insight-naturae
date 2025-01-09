// internal/mqttutils/mqttutils.go - project specific methods for handling MQTT communication
package mqttutils

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kajtekajtek/insight-naturae/pkg/utils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"strings"

	"github.com/joho/godotenv"
)

/* loads the MQTT connection options from the environment 
	variables and returns them as a MqttConn struct */
func LoadConnOpts() models.MqttConn {
	var c models.MqttConn

	// load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("Error loading .env file")
	}

	c.Scheme = utils.Getenv("MQTT_SCHEME", "tcp")
	c.Host = utils.Getenv("MQTT_HOST", "localhost")
	c.Port = utils.Getenv("MQTT_PORT", "1883")
	topics := utils.Getenv("MQTT_TOPICS", "insight-naturae/sensors")
	c.Topics = strings.Split(topics, ":")

	return c
}

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s: %s\n", msg.Topic(), msg.Payload())
}