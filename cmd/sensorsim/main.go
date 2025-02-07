// cmd/sensorsim/main.go - main file for simulation of sensors data
package main

import (
	"fmt"
	"log"

	"github.com/kajtekajtek/insight-naturae/internal/sensors"
	"github.com/kajtekajtek/insight-naturae/internal/config"
	"github.com/kajtekajtek/insight-naturae/pkg/mqtt"
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

	// load the configuratioon
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed while loading the configuration: %v", err)
	}

	// initialize the mqtt client and connect to the broker
	client, err := mqtt.InitClient(conf.MQTTScheme, conf.MQTTHost, conf.MQTTPort)
	if err != nil {
		log.Fatalf("Failed while initializing MQTT client: %v", err)
	}

	fmt.Println("Connected to MQTT broker on " + conf.MQTTScheme)
	for _, t := range conf.Topics {
		fmt.Println("Publishing on topic: " + t)
	}

	// simulate sensors
	go sensors.SimulateSensor(client, conf.Topics, u_temp, min_temp, max_temp, interval)
	go sensors.SimulateSensor(client, conf.Topics, u_hum, min_hum, max_hum, interval)
	go sensors.SimulateSensor(client, conf.Topics, u_pres, min_pres, max_pres, interval)
	go sensors.SimulateSensor(client, conf.Topics, u_co2, min_co2, max_co2, interval)

	// wait forever
	for {}
}	