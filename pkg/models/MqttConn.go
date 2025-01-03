/* pkg/models/MqttConn.go - defines the MQTT configuration struct 
	and methods */
package models

type MqttConn struct {
	Scheme string // "tcp", "ssl", or "ws"
	Host string // IP address or domain name
	Port string // port on which the broker listens
	Topics []string // topics to subscribe to
}