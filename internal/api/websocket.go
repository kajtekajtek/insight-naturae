/* internal/api/websocket.go - types, functions and methods for handling 
	WebSocket connections */
package api

import (
	"log"
	"net/http"
	"sync"
	"time"

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
type WebSocketClient struct {
	Username string
	Conn     *websocket.Conn
}

// client manager struct
type WebSocketClientManager struct {
	// clients map
	Clients map[string]*WebSocketClient // [username] => client struct
	// mutex for data protection
	Mutex   sync.RWMutex
}

// create a new client manager
func NewWebSocketClientManager() *WebSocketClientManager {
	return &WebSocketClientManager{
		Clients: make(map[string]*WebSocketClient),
	}
}

// add a new client to the manager
func (cm *WebSocketClientManager) AddClient(username string, conn *websocket.Conn) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Clients[username] = &WebSocketClient{Username: username, Conn: conn}
}

// remove a client from the manager
func (cm *WebSocketClientManager) RemoveClient(username string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	if client, found := cm.Clients[username]; found {
		if client.Conn != nil {
			client.Conn.Close()
		}
		delete(cm.Clients, username)
		log.Printf("Client %s disconnected", username)
	}
}

// broadcast a message to all clients
func (cm *WebSocketClientManager) Broadcast(message []byte) {
	cm.Mutex.RLock()
	defer cm.Mutex.RUnlock()
	// loop through all clients and send message
	for _, client := range cm.Clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// log error and close connection if failed to write message
			log.Printf("Failed to write message to client %s: %v", client.Username, err)
			client.Conn.Close()
		}
	}
}

// handle websocket connections
func (cm *WebSocketClientManager) WebSocketHandler(c *gin.Context) {
	// get the username from the header
	username := c.GetHeader("username")
	if username == ""  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
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

	// ping pong to keep connection alive
	// handle pong
	conn.SetPongHandler(func(appData string) error {
		log.Printf("Received pong from client %s", username)
		return nil
	})
	// send ping
	for {
		err := conn.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			log.Printf("Failed to send ping to client %s: %v", username, err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}