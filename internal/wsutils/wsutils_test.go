/* internal/wsutils/wsutils_test.go - tests for the WebSocket 	
	connection types, functions and methods */
package wsutils

import (
	"net/http/httptest"
	"time"
	"testing"
	"strconv"
	"database/sql"
	"os"
	"crypto/tls"

	"github.com/kajtekajtek/insight-naturae/internal/jwtutils"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	username = "testuser"
	JWTSecret = "testjwtsecretkey"
	DBPath = "test.db"
)

func SetupTestDB(t *testing.T) *sql.DB {
	// delete the test database if it exists
	if _, err := os.Stat(DBPath); err == nil {
		err := os.Remove(DBPath)
		assert.NoError(t, err)
	}

	// create the test database
	db, err := dbutils.CreateDatabase(DBPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	return db
}

func TearDownTestDB(db *sql.DB) {
	db.Close()
	os.Remove(DBPath)
}

func SetupWSServer(cm *WSClientManager, db *sql.DB) *httptest.Server {
	r := gin.Default()

	r.GET("/ws", cm.WebSocketHandler(db, []byte(JWTSecret)))

	testServer := httptest.NewUnstartedServer(r)
	testServer.StartTLS()

	return testServer
}

func SetupWSConnection(url string) (*websocket.Conn, error) {
	tlsConfig := &tls.Config{ InsecureSkipVerify: true }
	dialer := websocket.Dialer{TLSClientConfig: tlsConfig}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

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

func TestSubscribe(t *testing.T) {
	cm := NewWSClientManager()
	cm.Subscribe("testuser", "testsensor")

	assert.True(t, cm.Subscriptions["testsensor"]["testuser"])
}

func TestUnsubscribe(t *testing.T) {
	cm := NewWSClientManager()
	cm.Subscribe("testuser", "testsensor")
	cm.Unsubscribe("testuser", "testsensor")

	assert.False(t, cm.Subscriptions["testsensor"]["testuser"])
}

/* test the main WebSocket connection handler: create a client manager, 
	test server and connect to the server. Check if the client was added
		to the manager and removed after the connection was closed */
func TestWebSocketHandler(t *testing.T) {
	// setup the test database
	db := SetupTestDB(t)
	defer TearDownTestDB(db)	

	// create a new WSClientManager
	cm := NewWSClientManager()

	// create a test server
	t.Logf("Creating a test server...")
	server := SetupWSServer(cm, db)
	defer server.Close()
	t.Logf("Server URL: %s", server.URL)

	// create a JSON Web Token
	token, err := jwtutils.GenerateJWT([]byte(JWTSecret), username)
	assert.Nil(t, err)

	// connect to the server
	url := "wss" + server.URL[5:] + "/ws?token=" + token
	t.Logf("Connecting to %s", url)
	conn, err := SetupWSConnection(url)
	assert.Nil(t, err)
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

// test the WebSocket connection handler when no token is provided
func TestWebSocketHandlerNoToken(t *testing.T) {
	// setup the test database
	db := SetupTestDB(t)
	defer TearDownTestDB(db)	

	// create a new WSClientManager
	cm := NewWSClientManager()

	// create a test server
	server := SetupWSServer(cm, db)
	defer server.Close()

	// connect to the server without a token
	url := "wss" + server.URL[5:] + "/ws"
	_, err := SetupWSConnection(url)
	assert.NotNil(t, err)

	// check if the client was not added to the manager
	cm.Mutex.RLock()
	_, exists := cm.Clients[username]
	cm.Mutex.RUnlock()

	assert.False(t, exists)
}

/* test message broadcasting: create a client manager, test server and
	connect 3 clients to the server. Broadcast a message and check if
		each client received the message */
func TestBroadcast(t *testing.T) {
	// setup the test database
	db := SetupTestDB(t)
	defer TearDownTestDB(db)	

	// create a new WSClientManager
	cm := NewWSClientManager()
	
	// crate a test server
	server := SetupWSServer(cm, db)
	defer server.Close()

	url := "wss" + server.URL[5:] + "/ws" + "?token="

	// create 3 connections
	var conns []*websocket.Conn
	for i := 0; i < 3; i++ {
		// create a JSON Web Token
		token, err := jwtutils.GenerateJWT([]byte(JWTSecret), 
			username + strconv.Itoa(i))
		assert.Nil(t, err)
		
		// connect to the server
		t.Logf("Connecting to %s", url + token)
		conn, err := SetupWSConnection(url + token)
		assert.Nil(t, err)
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
	// setup the test database	
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	// create a new WSClientManager
	cm := NewWSClientManager()

	// create a test server
	server := SetupWSServer(cm, db)
	defer server.Close()

	// create a JSON Web Token
	token, err := jwtutils.GenerateJWT([]byte(JWTSecret), username)
	assert.Nil(t, err)

	// connect to the server
	url := "wss" + server.URL[5:] + "/ws" + "?token=" + token
	conn, err := SetupWSConnection(url)
	assert.Nil(t, err)
	defer conn.Close()

	// read the ping message first
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
	// setup the test database
	db := SetupTestDB(t)
	defer TearDownTestDB(db)	

	// create a new WSClientManager
	cm := NewWSClientManager()

	// create a test server
	server := SetupWSServer(cm, db)
	defer server.Close()

	// create a JSON Web Token
	token, err := jwtutils.GenerateJWT([]byte(JWTSecret), username)
	assert.Nil(t, err)

	// connect to the server
	url := "wss" + server.URL[5:] + "/ws" + "?token=" + token
	conn, err := SetupWSConnection(url)
	assert.Nil(t, err)
	conn.Close()

	// send a message
	cm.SendMessage(username, []byte("Hello, World!"))

	_, _, err = conn.ReadMessage()
	assert.NotNil(t, err)
}