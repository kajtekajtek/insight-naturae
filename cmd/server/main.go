// cmd/server/main.go - main entry point for the server application
package main

import (
	"log"

	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/pkg/mqtt"
	"github.com/kajtekajtek/insight-naturae/internal/mqttutils"
)

func main() {
	db, err := dbutils.CreateDatabase()
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	// load MQTT connection options
	conf := mqttutils.LoadConnOpts()

	// initialize the MQTT client
	mqttClient, err := mqtt.InitClient(conf.Scheme, conf.Host, conf.Port)
	if err != nil {
		log.Fatalf("Error initializing MQTT client: %v", err)
	}

	// subscribe to the topics
	messageHandler := mqttutils.MessageHandler(db)
	for _, t := range conf.Topics {
		log.Printf("Subscribing to topic: %s\n", t)
		if token := mqttClient.Subscribe(t, 0, messageHandler); token.Wait() && token.Error() != nil {
			log.Fatalf("Error subscribing to topic: %v", token.Error())
		}
	}

	// wait forever
	for {}
}