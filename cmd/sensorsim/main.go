// cmd/sensorsim/main.go - main file for simulation of sensors data
package main

import (
	"fmt"
	"log"

	"github.com/kajtekajtek/insight-naturae/internal/sensors"
	"github.com/kajtekajtek/insight-naturae/pkg/mqtt"
	"github.com/kajtekajtek/insight-naturae/pkg/utils"
)

func main() {
	// define the value ranges for the simulated sensors
	var min_temp, max_temp float64 = 10, 30
	var min_hum,  max_hum float64 = 0, 100
	var min_pres, max_pres float64 = 800, 1000
	var min_co2,  max_co2 float64  = 320, 520
	var u_temp, u_hum, u_pres, u_co2 string = "C", "%RH", "hPa", "PPM"
	// interval between sensor readings in seconds
	var interval int = 5

	// get the topic and broker address from environment variables
	scheme := utils.Getenv("MQTT_SCHEME", "tcp")
	host := utils.Getenv("MQTT_HOST", "localhost")
	port := utils.Getenv("MQTT_PORT", "1883")
	topic := utils.Getenv("MQTT_TOPIC", "insight-naturae/sensors")

	// initialize the mqtt client and connect to the broker
	client, err := mqtt.InitClient(scheme, host, port)
	if err != nil {
		log.Fatalf("Error initializing MQTT client: %v", err)
	}

	fmt.Println("Connected to MQTT broker on " + host)
	fmt.Println("Publishing to topic " + topic)

	// simulate sensors
	go sensors.SimulateSensor(client, topic, u_temp, min_temp, max_temp, interval)
	go sensors.SimulateSensor(client, topic, u_hum, min_hum, max_hum, interval)
	go sensors.SimulateSensor(client, topic, u_pres, min_pres, max_pres, interval)
	go sensors.SimulateSensor(client, topic, u_co2, min_co2, max_co2, interval)

	// wait forever
	for {}
}	