// cmd/server/main.go - main entry point for the server application
package main

import (
	"log"

	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/internal/mqttutils"
	"github.com/kajtekajtek/insight-naturae/internal/api"	
	"github.com/kajtekajtek/insight-naturae/internal/config"
	"github.com/kajtekajtek/insight-naturae/internal/middleware"
	"github.com/kajtekajtek/insight-naturae/pkg/mqtt"

	"github.com/gin-gonic/gin"
)

func main() {
	// load the configuration
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed while loading the configuration: %v", err)
	}

	// create the database
	db, err := dbutils.CreateDatabase(conf.DBPath)
	if err != nil {
		log.Fatalf("Failed while creating database: %v", err)
	}

	// initialize the MQTT client
	mqttClient, err := mqtt.InitClient(conf.MQTTScheme, conf.MQTTHost, conf.MQTTPort)
	if err != nil {
		log.Fatalf("Failed while initializing MQTT client: %v", err)
	}

	// subscribe to the topics
	messageHandler := mqttutils.MessageHandler(db)
	for _, t := range conf.Topics {
		log.Printf("Subscribing to topic: %s\n", t)
		if token := mqttClient.Subscribe(t, 0, messageHandler); token.Wait() && token.Error() != nil {
			log.Fatalf("Failed while subscribing to topic: %v", token.Error())
		}
	}

	// run the API server
	router := gin.Default()
	// public routes
	router.POST("/register", api.RegisterHandler(db))
	router.POST("/login", api.LoginHandler(db, conf.JWTSecret))
	// protected routes
	protected := router.Group("/user", middleware.AuthMiddleware(conf.JWTSecret))
	{
		protected.POST("/sensors", api.AddSensorHandler(db))
		protected.GET("/sensors", api.GetSensorsHandler(db))
		protected.DELETE("/sensors/:id", api.RemoveSensorHandler(db))
	}

	router.Run(":8080")
}