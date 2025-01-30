/* internal/websocket/websocket.go - types, functions and methods for handling 
	WebSocket connections */
package websocket

import (
	"log"
	"fmt"
	"net/http"
	"sync"
	"time"
	"database/sql"

	"github.com/kajtekajtek/insight-naturae/internal/dbutils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/* web socket upgrader - to upgrade the HTTP connection 
	to a web socket connection */
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all connections by default
	},
}

// client connection struct
type WSClient struct {
	Username string
	Conn     *websocket.Conn
	Mutex	 sync.Mutex
}

// client manager struct
type WSClientManager struct {
	// [username] => client struct
	Clients map[string]*WSClient
	// [sensorID] => [username] => subscribed
	Subscriptions map[string]map[string]bool 
	Mutex   sync.RWMutex
}

// create a new client manager
func NewWSClientManager() *WSClientManager {
	return &WSClientManager{
		Clients: make(map[string]*WSClient),
		Subscriptions: make(map[string]map[string]bool),
	}
}

// add a new client to the manager
func (cm *WSClientManager) AddClient(username string, conn *websocket.Conn) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Clients[username] = &WSClient{Username: username, Conn: conn, Mutex: sync.Mutex{}}
}

// remove a client from the manager
func (cm *WSClientManager) RemoveClient(username string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	if client, found := cm.Clients[username]; found {
		if client.Conn != nil {
			client.Conn.Close()
		}
		delete(cm.Clients, username)
		log.Printf("Client %s removed from Client Manager", username)
	}
}

// add the given username to the subscription list for the given sensor
func (cm *WSClientManager) Subscribe(username string, sensorID string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()

	// create a new subscription map for given sensor if it doesn't exist
	if _, exists := cm.Subscriptions[sensorID]; !exists {
		cm.Subscriptions[sensorID] = make(map[string]bool)
	}
	// add the username to the subscription list
	cm.Subscriptions[sensorID][username] = true
}

// remove the given username from the given sensor's subscription list
func (cm *WSClientManager) Unsubscribe(username string, sensorID string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()

	// check if the sensor exists in the subscriptions map
	if subscribers, exists := cm.Subscriptions[sensorID]; exists {
		// remove the username from the sensor's subscription list
		delete(subscribers, username)
		// delete the sensor from the subscriptions map if no subscribers left
		if len(subscribers) == 0 {
			delete(cm.Subscriptions, sensorID)
		}
	}
}

// get sensor's subscribers
func (cm *WSClientManager) GetSubscribers(sensorID string) []string {
	cm.Mutex.RLock()
	defer cm.Mutex.RUnlock()

	subscribers := make([]string, 0)
	if sensorSubscribers, exists := cm.Subscriptions[sensorID]; exists {
		for subscriber := range sensorSubscribers {
			subscribers = append(subscribers, subscriber)
		}
	}
	return subscribers
}

// broadcast a message to all clients
func (cm *WSClientManager) Broadcast(message []byte) error {
	// lock the clients map mutex
	cm.Mutex.RLock()
	defer cm.Mutex.RUnlock()

	// loop through all clients and send message
	for _, client := range cm.Clients {
		client.Mutex.Lock() // lock the client mutex
		defer client.Mutex.Unlock()

		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		// log error and close connection if WriteMessage fails
		if err != nil {
			client.Conn.Close()

			log.Printf("Failed to write message to client %s: %v", 
				client.Username, err)
			return fmt.Errorf("failed to write message to client %s: %v", 
				client.Username, err)
		}
	}
	return nil
}

// send a message to a specific client
func (cm *WSClientManager) SendMessage(username string, message []byte) error {
	cm.Mutex.RLock()
	client, found := cm.Clients[username]
	cm.Mutex.RUnlock()

	if found {
		client.Mutex.Lock()
		defer client.Mutex.Unlock()
		err := client.Conn.WriteMessage(websocket.TextMessage, message); 
		// log error and close connection if failed to write message
		if err != nil {
			client.Conn.Close()

			log.Printf("Failed to write message to client %s: %v", 
				username, err)
			return fmt.Errorf("failed to write message to client %s: %v", 
				username, err)
		}
	}
	return nil
}

/* ping pong function keeps the connection alive by sending ping messages 
	and receiving pong messages, updating the connection status channel */
func PingPong(connected chan bool, cm *WSClientManager, username string) {
	cm.Mutex.RLock()
	client, found := cm.Clients[username]
	if !found {
		connected <- false
		return
	}
	conn := client.Conn
	cm.Mutex.RUnlock()

	// handle pong
	conn.SetPongHandler(func(appData string) error {
		log.Printf("Received pong from client %s", username)
		return nil
	})
	// send ping
	for {
		err := cm.SendMessage(username, []byte("ping"))
		if err != nil {
			log.Printf("Failed to send ping to client %s: %v", 
				username, err)
			connected <- false
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// handle websocket connections
func (cm *WSClientManager) WebSocketHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the username from the header
		username := c.GetHeader("username")
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No username provided"})
			return
		}

		// upgrade the HTTP connection to a websocket connection
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection to websocket: %v", err)
			return
		}

		// add the client to the manager
		cm.AddClient(username, conn)
		defer cm.RemoveClient(username)

		log.Printf("WebSocket connection with client %s established", username)

		// get user's subscriptions from the database
		sensors, err := dbutils.GetUserSubscriptions(db, username);
		if err != nil {
			log.Printf("Failed to get user %s subscriptions: %v", 
				username, err)
		} else {
			// add the user to the subscriptions map for each sensor
			for _, sensor := range sensors {
				cm.Subscribe(username, sensor.SensorID)
			}
		}

		// ping pong to keep the connection alive
		status := make(chan bool)
		go PingPong(status, cm, username)

		for {
			select {
			// get the connection status from the ping pong function
			case connected := <-status:
				if !connected {
					return
				}
			}
		}
	}
}