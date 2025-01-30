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
	ws "github.com/kajtekajtek/insight-naturae/internal/wsutils"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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

	// create the WebSocket client manager
	clientManager := ws.NewWSClientManager()

	// subscribe to the topics
	messageHandler := mqttutils.MessageHandler(db, clientManager)
	for _, t := range conf.Topics {
		log.Printf("Subscribing to topic: %s\n", t)
		if token := mqttClient.Subscribe(t, 0, messageHandler); token.Wait() && token.Error() != nil {
			log.Fatalf("Failed while subscribing to topic: %v", token.Error())
		}
	}

	// set up the server
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: c.CORSOrigins,
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// public routes
	router.POST("/register", api.RegisterHandler(db))
	router.POST("/login", api.LoginHandler(db, conf.JWTSecret))
	router.GET("/ws", clientManager.WebSocketHandler(db, conf.JWTSecret))

	// protected routes
	protected := router.Group("/user", middleware.AuthMiddleware(conf.JWTSecret))
	{
		protected.POST("/sensors", api.SubscribeSensorHandler(db))
		protected.GET("/sensors", api.GetSensorsHandler(db))
		protected.DELETE("/sensors/:id", api.UnsubscribeSensorHandler(db))
	}

	router.Run(":8080")
}