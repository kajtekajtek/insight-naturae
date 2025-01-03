// internal/mqttutils/mqttutils.go - package for handling MQTT communication
package mqttutils

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kajtekajtek/insight-naturae/pkg/utils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"strings"
)

func LoadConnOpts() models.MqttConn {
	var c models.MqttConn
	c.Scheme = utils.Getenv("MQTT_SCHEME", "tcp")
	c.Host = utils.Getenv("MQTT_HOST", "localhost")
	c.Port = utils.Getenv("MQTT_PORT", "1883")
	topics := utils.Getenv("MQTT_TOPICS", "insight-naturae/sensors")
	c.Topics = strings.Split(topics, ":")

	return c
}