// cmd/sensorsim/main.go - main file for simulation of sensors data
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/kajtekajtek/insight-naturae/pkg/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kajtekajtek/insight-naturae/internal/sensors"
	"github.com/google/uuid" // uuid
)

func main() {
	// default mqtt broker options
	var defaultURL string = "localhost:1883"
	var defaultTopic string = "insight-naturae/sensors"
	// define the value ranges for the simulated sensors
	var min_temp, max_temp float64 = 10, 30
	var min_hum,  max_hum float64 = 0, 100
	var min_pres, max_pres float64 = 800, 1000
	var min_co2,  max_co2 float64  = 320, 520
	var u_temp, u_hum, u_pres, u_co2 string = "C", "%RH", "hPa", "PPM"
	// interval between sensor readings in seconds
	var interval int = 5

	// get the topic and URL from environment variables
	topic := utils.Getenv("MQTT_TOPIC", defaultTopic)
	url := utils.Getenv("MQTT_URL", defaultURL)

	// generate an uuid, encode it to base64 and set it as the client ID
	id := uuid.New()
	clientID := base64.RawURLEncoding.EncodeToString(id[:])

	// connect to the MQTT broker
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + url).SetClientID(clientID)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Connected to MQTT broker on " + url)
	fmt.Println("Publishing to topic " + topic)
	fmt.Println("Client ID: " + opts.ClientID)

	// simulate sensors
	go sensors.SimulateSensor(client, topic, u_temp, min_temp, max_temp, interval)
	go sensors.SimulateSensor(client, topic, u_hum, min_hum, max_hum, interval)
	go sensors.SimulateSensor(client, topic, u_pres, min_pres, max_pres, interval)
	go sensors.SimulateSensor(client, topic, u_co2, min_co2, max_co2, interval)

	// wait forever
	for {}
}	