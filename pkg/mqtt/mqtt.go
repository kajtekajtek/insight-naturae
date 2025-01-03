// pkg/mqtt/mqtt.go - generic MQTT communication methods
package mqtt

import (
	"encoding/base64"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// InitClient initializes a new MQTT client and connects to the broker
func InitClient(scheme, host, port string) (mqtt.Client, error) {
	broker := scheme + "://" + host + ":" + port

	id := uuid.New()
	clientID := base64.RawURLEncoding.EncodeToString(id[:])

	// set the client options
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	// create a new client
	client := mqtt.NewClient(opts)
	// connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func Subscribe(client mqtt.Client, topic string, callback mqtt.MessageHandler) error {
	if token := client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}