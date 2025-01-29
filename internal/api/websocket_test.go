/* internal/api/websocket_test.go - tests for the WebSocket 	
	connection types, functions and methods */
package api

import (
	"net/http/httptest"
	"time"
	"net/http"
	"testing"
	"strconv"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestNewWSClientManager(t *testing.T) {
	cm := NewWSClientManager()

	assert.NotNil(t, cm)
	assert.NotNil(t, cm.Clients)
}

func TestAddClient(t *testing.T) {
	cm := NewWSClientManager()
	cm.AddClient("testuser", nil)

	assert.NotNil(t, cm.Clients["testuser"])
}

func TestRemoveClient(t *testing.T) {
	cm := NewWSClientManager()
	cm.AddClient("testuser", nil)
	cm.RemoveClient("testuser")

	assert.Nil(t, cm.Clients["testuser"])
}

/* test the main WebSocket connection handler: create a client manager, 
	test server and connect to the server. Check if the client was added
		to the manager and removed after the connection was closed */
func TestWebSocketHandler(t *testing.T) {
	// create a new WSClientManager
	cm := NewWSClientManager()

	// create a test server
	server := SetupWSServer(t, cm)
	defer server.Close()

	url := "ws" + server.URL[4:] + "/ws"

	/* normally, the user should pass a JSON Web Token in the headers 
		and username would be extracted from it by the 
			authentification middleware */
	headers := http.Header{}
	headers.Set("username", username)

	// connect to the server
	conn := SetupWSConnection(t, url, headers)
	defer conn.Close()

	// check if the client was added to the manager
	cm.Mutex.RLock()
	_, exists := cm.Clients[username]
	cm.Mutex.RUnlock()

	assert.True(t, exists)

	// close the connection and check if the client was removed
	conn.Close()
	time.Sleep(2 * time.Second)

	cm.Mutex.RLock()
	_, exists = cm.Clients[username]
	cm.Mutex.RUnlock()

	assert.False(t, exists)
}

// test the WebSocket connection handler when no username is provided
func TestWebSocketHandlerNoUsername(t *testing.T) {
	// only setup the router
	r := gin.Default()
	r.GET("/ws", NewWSClientManager().WebSocketHandler)

	// create a request without an username provided
	req, _ := http.NewRequest("GET", "/ws", nil)

	// send the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No username provided")
}

/* test message broadcasting: create a client manager, test server and
	connect 3 clients to the server. Broadcast a message and check if
		each client received the message */
func TestBroadcast(t *testing.T) {
	// create a new WSClientManager
	cm := NewWSClientManager()
	
	// crate a test server
	server := SetupWSServer(t, cm)
	defer server.Close()

	url := "ws" + server.URL[4:] + "/ws"

	// create 3 connections
	var conns []*websocket.Conn
	for i := 0; i < 3; i++ {
		headers := http.Header{}
		headers.Set("username", username + strconv.Itoa(i)) 

		conn := SetupWSConnection(t, url, headers)
		conns = append(conns, conn)
	}
	defer func() {
		for _, conn := range conns { conn.Close() }
	}()

	// read the ping messages first
	for _, conn := range conns {
		_, msg, err := conn.ReadMessage()
		assert.Nil(t, err)
		assert.Equal(t, "ping", string(msg))
	}

	// broadcast a message
	cm.Broadcast([]byte("Hello, World!"))

	// check if each client received the message
	for _, conn := range conns {
		_, msg, err := conn.ReadMessage()
		assert.Nil(t, err)
		assert.Equal(t, "Hello, World!", string(msg))
	}
}

// test sending a message to a specific client
func TestSendMessage(t *testing.T) {
	// create a new WSClientManager
	cm := NewWSClientManager()

	// create a test server
	server := SetupWSServer(t, cm)
	defer server.Close()

	url := "ws" + server.URL[4:] + "/ws"

	// create a connection
	headers := http.Header{}
	headers.Set("username", username)

	conn := SetupWSConnection(t, url, headers)
	defer conn.Close()

	_, msg, err := conn.ReadMessage()
	assert.Nil(t, err)
	assert.Equal(t, "ping", string(msg))	

	// send a message
	cm.SendMessage(username, []byte("Hello, World!"))

	// check if the client received the message
	_, msg, err = conn.ReadMessage()
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(msg))
}

func TestSendMessageNoClient(t *testing.T) {
	// create a new WSClientManager
	cm := NewWSClientManager()

	// send a message to a non-existing client
	cm.SendMessage("nonexisting", []byte("Hello, World!"))

	// no error should be thrown
}

func TestSendMessageClosedConnection(t *testing.T) {
	// create a new WSClientManager
	cm := NewWSClientManager()

	// create a test server
	server := SetupWSServer(t, cm)
	defer server.Close()

	url := "ws" + server.URL[4:] + "/ws"

	// create a connection
	headers := http.Header{}
	headers.Set("username", username)

	conn := SetupWSConnection(t, url, headers)
	conn.Close()

	// send a message
	cm.SendMessage(username, []byte("Hello, World!"))

	_, _, err := conn.ReadMessage()
	assert.NotNil(t, err)
}