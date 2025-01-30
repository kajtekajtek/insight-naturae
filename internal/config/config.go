// internal/config/config.go - project configuration struct and methods
package config

import (
	"fmt"
	"strings"

	"github.com/kajtekajtek/insight-naturae/pkg/utils"

	"github.com/joho/godotenv"
)

type Config struct {
	MQTTScheme string // "tcp", "ssl", or "ws"
	MQTTHost string // IP address or domain name
	MQTTPort string // port on which the broker listens
	Topics []string // topics to subscribe to
	DBPath string // path to the SQLite database file
	JWTSecret []byte // secret key for JWT signing
	CORSOrigins []string // list of allowed origins for CORS
}

/* loads the configuration from the .env file if provided, 
	otherwise uses the default values */
func LoadConfig() (Config, error) {
	var c Config

	// load the environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file; continuing with the default values...")
	}

	// MQTT connection parameters
	c.MQTTScheme = utils.Getenv("MQTT_SCHEME", "tcp")
	c.MQTTHost = utils.Getenv("MQTT_HOST", "localhost")
	c.MQTTPort = utils.Getenv("MQTT_PORT", "1883")
	topics := utils.Getenv("MQTT_TOPICS", "insight-naturae/sensors")
	c.Topics = strings.Split(topics, ":")

	// Database path
	c.DBPath = utils.Getenv("DB_PATH", "insight-naturae.db")

	// JWT secret key used for signing the authentication tokens
	JWTSecretString := utils.Getenv("JWT_SECRET", "")	
	// if secret key is not provided, generate a new one
	if JWTSecretString == "" {
		var err error
		c.JWTSecret, err = utils.GenerateSecretKey(32)
		if err != nil {
			return Config{}, err
		}
	// otherwise, use the provided secret key
	} else {
		c.JWTSecret = []byte(JWTSecretString)
	}

	// CORS allowed origins
	origins := utils.Getenv("CORS_ORIGINS", "http://localhost:5000")
	c.CORSOrigins = strings.Split(origins, ",")

	return c, nil
}