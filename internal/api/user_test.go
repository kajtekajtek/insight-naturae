// internal/api/user_test.go - integration tests for the login & register endpoints
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
)

func TestRegisterHandler(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// test the register endpoint
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

// try to register existing user
func TestRegisterHandlerExistingUser(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// register the user
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// register the same user again
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "User already exists")
}

// try to register with password field missing
func TestRegisterHandlerFieldMissing(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
	}
	payload, _ := json.Marshal(user)

	// test the register endpoint
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// try to register with empty strings
func TestRegisterHandlerEmptyStrings(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: "",
		Password: "",
	}
	payload, _ := json.Marshal(user)

	// test the register endpoint
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestLoginHandler(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// register the user
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// test the login endpoint
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

// try to login with invalid input
func TestLoginHandlerFieldMissing(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
	}
	payload, _ := json.Marshal(user)

	// test the login endpoint
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// try to login to a non-existing user
func TestLoginHandlerNonExistingUser(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// test the login endpoint
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}

// try to login with invalid password
func testLoginHandlerInvalidCredentials(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: username,
		Password: password,
	}
	payload, _ := json.Marshal(user)

	// register the user
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// create the user payload
	user = models.User{
		Username: username,
		Password: "wrongpassword",
	}
	payload, _ = json.Marshal(user)

	// test the login endpoint
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}

// try to login with empty strings
func testLoginHandlerEmptyStrings(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	r := SetupRouter(db)

	// create the user payload
	user := models.User{
		Username: "",
		Password: "",
	}
	payload, _ := json.Marshal(user)

	// test the login endpoint
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}