// internal/api/user_test.go - integration tests for the login & register endpoints
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
)

const testDBPath = "test.db"
const JWTSecret = "testjwtsecretkey"
const username = "testuser"
const password = "testpassword"

func SetupTestDB(t *testing.T) *sql.DB {
	// delete the test database if it exists
	if _, err := os.Stat(testDBPath); err == nil {
		err := os.Remove(testDBPath)
		assert.NoError(t, err)
	}

	// create the test database
	db, err := dbutils.CreateDatabase(testDBPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	return db
}

func TearDownTestDB(db *sql.DB) {
	db.Close()
	os.Remove(testDBPath)
}

func setupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/register", RegisterHandler(db))
	r.POST("/login", LoginHandler(db, []byte(JWTSecret)))

	return r
}

func TestRegisterHandler(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := setupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// test the register endpoint
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

func TestRegisterHandlerConflict(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := setupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// register the user
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// register the same user again
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "User already exists")
}

func TestRegisterHandlerInvalidInput(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := setupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
	}
	payload, _ := json.Marshal(user)

	// test the register endpoint
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestLoginHandler(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := setupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// register the user
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// test the login endpoint
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginHandlerInvalidInput(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := setupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
	}
	payload, _ := json.Marshal(user)

	// test the login endpoint
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := setupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// test the login endpoint
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}